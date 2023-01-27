package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (c *Client) terminatePersistentQuery(ctx context.Context, name string) error {
	var (
		err error
		res []QueryResponse
	)

	err = c.makePostKsqlRequestWithUnmarshal(ctx, "SHOW QUERIES;",
		func(r io.Reader) error { return json.NewDecoder(r).Decode(&res) },
	)
	if err != nil {
		return err
	}

	terminateQueries := make([]string, 0, 1)

	for _, query := range res {
		for _, pq := range query.Queries {
			for _, s := range pq.Sinks {
				if strings.ToLower(s) == strings.ToLower(name) {
					terminateQueries = append(terminateQueries, pq.ID)
				}
			}
		}
	}

	var termRes []map[string]interface{}

	err = c.makePostKsqlRequestWithUnmarshal(ctx,
		fmt.Sprintf("TERMINATE %s ;", strings.Join(terminateQueries, ", ")),
		func(r io.Reader) error { return json.NewDecoder(r).Decode(&termRes) },
	)
	if err != nil {
		return err
	}

	return nil
}
