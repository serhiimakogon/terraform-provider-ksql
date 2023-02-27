package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	cli                 *http.Client
	url                 string
	username            string
	password            string
	autoOffsetResetMode string
}

func New(url, username, password, autoOffsetResetMode string) *Client {
	return &Client{
		url:                 url,
		username:            username,
		password:            password,
		autoOffsetResetMode: autoOffsetResetMode,
		cli:                 &http.Client{},
	}
}

func (c *Client) RotateCredentials(url, username, password string) {
	if url != "" {
		c.url = url
	}
	if username != "" {
		c.username = username
	}
	if password != "" {
		c.password = password
	}
}

func (c *Client) ResolveAutoOffsetResetQueryProperty(mode string) string {
	const pattern = "SET 'auto.offset.reset'='%s';"

	if mode != "" && mode != "0" {
		return fmt.Sprintf(pattern, mode)
	}
	if c.autoOffsetResetMode != "" {
		return fmt.Sprintf(pattern, c.autoOffsetResetMode)
	}

	return ""
}

func (c *Client) ExecuteQuery(ctx context.Context, name, qType, query string, ignoreAlreadyExists, terminatePersistentQuery bool) (string, error) {
	var (
		err           error
		res           interface{}
		resErrCode    float64
		resErrMessage string
	)

	for _, backoff := range []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		6 * time.Second,
		8 * time.Second,
		10 * time.Second,
	} {
		err = c.makePostKsqlRequestWithUnmarshal(ctx, query,
			func(r io.Reader) error { return json.NewDecoder(r).Decode(&res) },
		)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("failed to make post ksql request [%v] retrying...", err))
			time.Sleep(backoff)
			continue
		}

		switch g := res.(type) {
		case []interface{}:
			if len(g) == 0 {
				break
			}

			resErrCode, _ = g[0].(map[string]interface{})["error_code"].(float64)
			resErrMessage, _ = g[0].(map[string]interface{})["message"].(string)
		case map[string]interface{}:
			resErrCode, _ = g["error_code"].(float64)
			resErrMessage, _ = g["message"].(string)
		}

		if resErrCode != 0 {
			if ignoreAlreadyExists && strings.Contains(resErrMessage, "already exists") {
				break
			}
			if terminatePersistentQuery ||
				strings.Contains(resErrMessage, "Upgrades not yet supported") ||
				strings.Contains(resErrMessage, "Cannot drop") {
				if err = c.terminatePersistentQuery(ctx, name); err != nil {
					err = fmt.Errorf("failed to terminate persistent query: %v", err)
				}
				continue
			}

			err = fmt.Errorf("invalid ksql response %s", resErrMessage)
			if strings.HasPrefix(query, "DROP") {
				if terminateQuery, shouldTerminate := c.getPreHookTerminateQuery(resErrMessage); shouldTerminate {
					_, err = c.ExecuteQuery(ctx, name, qType, terminateQuery, ignoreAlreadyExists, terminatePersistentQuery)
				}
			}
			tflog.Warn(ctx, fmt.Sprintf("failed to make post ksql request [%v] retrying...", err))
			time.Sleep(backoff)
			continue
		}

		break
	}

	if err != nil {
		return "", err
	}

	return qType + "_" + name, nil
}

func (c *Client) getPreHookTerminateQuery(msg string) (string, bool) {
	queries := make([]string, 0)
	for _, line := range strings.Split(msg, "\n") {
		if strings.HasPrefix(line, "The following queries") {
			if items := strings.TrimSpace(line[strings.Index(line, "[")+1 : strings.Index(line, "]")]); items != "" {
				queries = append(queries, strings.Split(items, ",")...)
			}
		}
	}
	if len(queries) == 0 {
		return "", false
	}

	return "TERMINATE " + strings.Join(queries, ", ") + ";", true
}

func (c *Client) makePostKsqlRequestWithUnmarshal(ctx context.Context, query string, unmarshal func(in io.Reader) error) error {
	b, err := json.Marshal(map[string]interface{}{"ksql": query})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/ksql", c.url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = unmarshal(resp.Body); err != nil {
		return err
	}

	return nil
}

func ExtractNameFromQuery(query string) string {
	// keywords that may be present before object name
	var keywords = map[string]bool{
		"CREATE":  true,
		"DROP":    true,
		"OR":      true,
		"REPLACE": true,
		"TABLE":   true,
		"STREAM":  true,
	}

	words := strings.Split(query, " ")
	for _, word := range words {
		if !keywords[word] {
			return word
		}
	}

	return uuid.New().String() // random
}
