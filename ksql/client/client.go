package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
	"time"
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
		err error
		res Response
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
		if err == nil {
			break
		}

		tflog.Warn(ctx, fmt.Sprintf("failed to make post ksql request [%v] retrying...", err))
		time.Sleep(backoff)
	}

	for i, r := range res {
		if r.ErrorCode != 0 {
			return "", errors.New(res[i].Message)
		}
	}

	return qType + "_" + name, nil
}

func (c *Client) makePostKsqlRequest(ctx context.Context, query string) (Response, error) {
	b, err := json.Marshal(map[string]interface{}{"ksql": query})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/ksql", c.url), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if sc := resp.StatusCode; sc < 200 || sc > 300 {
		return nil, fmt.Errorf("invalid response status code [%d], body [%s]", sc, string(body))
	}

	res := Response{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s, err: %v", string(body), err)
	}

	return res, nil
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
