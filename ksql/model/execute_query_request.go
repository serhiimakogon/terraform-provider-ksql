package model

import (
	"bytes"
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ExecuteQueryRequest struct {
	name                     string
	queryType                string
	queryContent             string
	resourceType             string
	ignoreAlreadyExists      bool
	deleteTopicOnDestroy     bool
	terminatePersistentQuery bool
	properties               *QueryProperties
}

func NewExecuteQueryRequest(
	name, queryType, queryContent, resourceType string,
	ignoreAlreadyExists, deleteTopicOnDestroy, terminatePersistentQuery bool,
	properties *QueryProperties) *ExecuteQueryRequest {
	return &ExecuteQueryRequest{
		name:                     name,
		queryType:                queryType,
		queryContent:             sanitizeQueryContent(queryContent),
		resourceType:             resourceType,
		ignoreAlreadyExists:      ignoreAlreadyExists,
		deleteTopicOnDestroy:     deleteTopicOnDestroy,
		terminatePersistentQuery: terminatePersistentQuery,
		properties:               properties.MergeWithQueryContent(queryContent),
	}
}

func (eqr *ExecuteQueryRequest) ID() string {
	return eqr.resourceType + "_" + eqr.name
}

func (eqr *ExecuteQueryRequest) ResourceName() string {
	return eqr.name
}

func (eqr *ExecuteQueryRequest) IsDestroy() bool {
	return eqr.queryType == "destroy"
}

func (eqr *ExecuteQueryRequest) CheckAlreadyExistsError(errMessage string) bool {
	return eqr.ignoreAlreadyExists && strings.Contains(errMessage, "already exists")
}

func (eqr *ExecuteQueryRequest) CheckPersistentQueryDependencyError(errMessage string) bool {
	return eqr.terminatePersistentQuery ||
		strings.Contains(errMessage, "Upgrades not yet supported") ||
		strings.Contains(errMessage, "Cannot drop")
}

func (eqr *ExecuteQueryRequest) GenerateQueryContent() string {
	ctx := context.Background()

	switch eqr.queryType {
	case "create":
		return eqr.generateCreateQuery()
	case "delete":
		return eqr.generateDeleteQuery()
	case "read":
		tflog.Info(ctx, "do nothing", map[string]interface{}{"query_type": eqr.queryType})
		return ""
	default:
		tflog.Error(ctx, "invalid query type", map[string]interface{}{"query_type": eqr.queryType})
		return ""
	}
}

func (eqr *ExecuteQueryRequest) generateCreateQuery() string {
	return eqr.properties.ToQueryContent() + eqr.queryContent
}

func (eqr *ExecuteQueryRequest) generateDeleteQuery() string {
	buf := &bytes.Buffer{}
	buf.WriteString("DROP ")
	buf.WriteString(eqr.resourceType)
	buf.WriteString(" IF EXISTS ")
	buf.WriteString(eqr.name)
	if eqr.deleteTopicOnDestroy {
		buf.WriteString(" DELETE TOPIC ")
	}
	buf.WriteString(" ;")

	return buf.String()
}

func sanitizeQueryContent(content string) string {
	parts := strings.Split(content, ";")
	res := make([]string, 0, len(parts))

	for _, part := range parts {
		if strings.HasPrefix(part, "SET") {
			continue
		}

		res = append(res, part)
	}

	return strings.Join(res, " ")
}
