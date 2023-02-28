package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"terraform-provider-ksql/ksql/model"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	cli        *http.Client
	url        string
	username   string
	password   string
	properties *model.QueryProperties
}

func New(url, username, password string, properties *model.QueryProperties) *Client {
	return &Client{
		url:        url,
		username:   username,
		password:   password,
		properties: properties,
		cli:        &http.Client{},
	}
}

func (c *Client) MergeWithGlobalProperties(in model.QueryProperties) model.QueryProperties {
	return c.properties.Merge(in)
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

func (c *Client) ExecuteQuery(ctx context.Context, query *model.ExecuteQueryRequest) error {
	var (
		res any
		err error
	)

	for _, backoff := range []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		6 * time.Second,
		8 * time.Second,
		10 * time.Second,
	} {
		err = c.makePostKsqlRequestWithUnmarshal(ctx, query.GenerateQueryContent(),
			func(r io.Reader) error { return json.NewDecoder(r).Decode(&res) },
		)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("failed to make post ksql request [%v] retrying...", err))
			time.Sleep(backoff)
			continue
		}

		errCode, errMessage := c.parseErrorResponse(res)
		if errCode == 0 {
			break
		}

		if query.CheckAlreadyExistsError(errMessage) {
			break
		}

		if query.CheckPersistentQueryDependencyError(errMessage) {
			if err := c.terminatePersistentQuery(ctx, query.ResourceName()); err != nil {
				tflog.Error(ctx, "terminate persistent query", map[string]interface{}{"err": err})
			}
		}

		tflog.Warn(ctx, "make post ksql request retrying...", map[string]interface{}{"err": err})
		time.Sleep(backoff)
	}

	return err
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

func (c *Client) parseErrorResponse(res any) (float64, string) {
	switch g := res.(type) {
	case []interface{}:
		if len(g) == 0 {
			return 0, ""
		}
		errCode, _ := g[0].(map[string]interface{})["error_code"].(float64)
		errMessage, _ := g[0].(map[string]interface{})["message"].(string)

		return errCode, errMessage

	case map[string]interface{}:
		errCode, _ := g["error_code"].(float64)
		errMessage, _ := g["message"].(string)

		return errCode, errMessage

	default:
		return 0, ""
	}
}
