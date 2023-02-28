package model

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type QueryProperties struct {
	val map[string]string
}

func NewQueryProperties(in map[string]interface{}) QueryProperties {
	qp := QueryProperties{val: make(map[string]string, len(in))}

	for key, value := range in {
		val, ok := value.(string)
		if !ok {
			tflog.Warn(context.Background(), "invalid global query property value",
				map[string]interface{}{"key": key, "value": value})
			continue
		}

		qp.val[key] = val
	}

	return qp
}

func (qp QueryProperties) Merge(in QueryProperties) QueryProperties {
	res := QueryProperties{val: make(map[string]string)}

	for key, val := range qp.val {
		res.val[key] = val
	}

	for key, val := range in.val {
		res.val[key] = val
	}

	return res
}

func (qp QueryProperties) MergeWithQueryContent(content string) QueryProperties {
	res := QueryProperties{val: make(map[string]string)}

	for key, val := range qp.val {
		res.val[key] = val
	}

	for _, part := range strings.Split(content, ";") {
		if strings.HasPrefix(part, "SET") {
			part = strings.TrimPrefix(part, "SET ")
			kvParts := strings.Split(part, "=")
			key := strings.TrimLeft(strings.TrimRight(kvParts[0], "'"), "'")
			val := strings.TrimLeft(strings.TrimRight(kvParts[1], "'"), "'")

			res.val[key] = val
		}
	}

	return res
}

func (qp QueryProperties) ToQueryContent() string {
	buf := &bytes.Buffer{}

	for key, val := range qp.val {
		buf.WriteString(fmt.Sprintf("SET '%s'='%s';", key, val))
	}

	return buf.String()
}
