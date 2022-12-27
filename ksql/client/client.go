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
	cli      *http.Client
	url      string
	username string
	password string
}

func New(url, username, password string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
		cli:      &http.Client{},
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

func (c *Client) ExecuteQuery(ctx context.Context, name, qType, query string) (string, error) {
	var (
		queryModified bool
		err           error
		res           Response
	)

	for _, backoff := range []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		6 * time.Second,
		8 * time.Second,
		10 * time.Second,
	} {
		res, err = c.makePostKsqlRequest(ctx, query)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("failed to make post ksql request [%v] retrying...", err))
			time.Sleep(backoff)
			continue
		}

		if res.ErrorCode != 0 {
			err = fmt.Errorf("invalid ksql response %s", res.Message)
			if strings.HasPrefix(query, "DROP") && queryModified {
				if terminateQuery, shouldModify := c.getPreHookTerminateQuery(res.Message); shouldModify {
					query = terminateQuery + " " + query
				}
				queryModified = true
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

	return "TERMINATE " + strings.Join(queries, ", ") + " ;", true
}

func (c *Client) makePostKsqlRequest(ctx context.Context, query string) (Response, error) {
	b, err := json.Marshal(map[string]interface{}{"ksql": query})
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/ksql", c.url), bytes.NewBuffer(b))
	if err != nil {
		return Response{}, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.cli.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return Response{}, nil
	}

	res := &Response{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal: %s, err: %v", string(body), err)
	}

	return *res, nil
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
