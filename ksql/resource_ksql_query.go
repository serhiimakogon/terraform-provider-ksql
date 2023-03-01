package ksql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-ksql/ksql/client"
	"terraform-provider-ksql/ksql/model"
)

func resourceQuery() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a KSQL Query resource.",
		CreateContext: resourceQueryCreate,
		ReadContext:   resourceQueryRead,
		DeleteContext: resourceQueryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "KSQL query name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "KSQL query type [table|stream].",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"query": {
				Description: "KSQL query.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"delete_topic_on_destroy": {
				Description: "Delete topic on destroy.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"ignore_already_exists": {
				Description: "Ignore already exists errors.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"terminate_persistent_query": {
				Description: "Terminate persistent query if needed.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"query_properties": {
				Description: "Map of query properties",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Default:     map[string]interface{}{},
			},
			"credentials": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The KSQL Cluster API Credentials.",
				MinItems:    1,
				MaxItems:    1,
				Sensitive:   true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The KSQL URL.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The KSQL username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The KSQL password.",
						},
					},
				},
			},
		},
	}
}

func extractStringValueFromBlock(d *schema.ResourceData, blockName string, attribute string) string {
	v, ok := d.Get(fmt.Sprintf("%s.0.%s", blockName, attribute)).(string)
	if !ok {
		return ""
	}
	return v
}

func resourceQueryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cli := m.(*client.Client)
	cli.RotateCredentials(
		extractStringValueFromBlock(d, "credentials", "url"),
		extractStringValueFromBlock(d, "credentials", "username"),
		extractStringValueFromBlock(d, "credentials", "password"),
	)

	var (
		diags                    diag.Diagnostics
		queryType                = "create"
		name                     = d.Get("name").(string)
		queryContent             = d.Get("query").(string)
		resourceType             = d.Get("type").(string)
		ignoreAlreadyExists      = d.Get("ignore_already_exists").(bool)
		deleteTopicOnDestroy     = d.Get("delete_topic_on_destroy").(bool)
		terminatePersistentQuery = d.Get("terminate_persistent_query").(bool)
		queryProperties          = d.Get("query_properties").(map[string]interface{})
	)

	eqp := model.NewExecuteQueryRequest(
		name, queryType, queryContent, resourceType,
		ignoreAlreadyExists,
		deleteTopicOnDestroy,
		terminatePersistentQuery,
		cli.MergeWithGlobalProperties(model.NewQueryProperties(queryProperties)),
	)

	err := cli.ExecuteQuery(ctx, eqp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(eqp.ID())

	return diags
}

func resourceQueryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cli := m.(*client.Client)
	cli.RotateCredentials(
		extractStringValueFromBlock(d, "credentials", "url"),
		extractStringValueFromBlock(d, "credentials", "username"),
		extractStringValueFromBlock(d, "credentials", "password"),
	)

	var (
		diags                    diag.Diagnostics
		queryType                = "read"
		name                     = d.Get("name").(string)
		queryContent             = d.Get("query").(string)
		resourceType             = d.Get("type").(string)
		ignoreAlreadyExists      = d.Get("ignore_already_exists").(bool)
		deleteTopicOnDestroy     = d.Get("delete_topic_on_destroy").(bool)
		terminatePersistentQuery = d.Get("terminate_persistent_query").(bool)
		queryProperties          = d.Get("query_properties").(map[string]interface{})
	)

	eqp := model.NewExecuteQueryRequest(
		name, queryType, queryContent, resourceType,
		ignoreAlreadyExists,
		deleteTopicOnDestroy,
		terminatePersistentQuery,
		cli.MergeWithGlobalProperties(model.NewQueryProperties(queryProperties)),
	)

	tflog.Info(ctx, fmt.Sprintf("Going to read object [%s]", eqp.ID()))

	return diags
}

func resourceQueryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cli := m.(*client.Client)
	cli.RotateCredentials(
		extractStringValueFromBlock(d, "credentials", "url"),
		extractStringValueFromBlock(d, "credentials", "username"),
		extractStringValueFromBlock(d, "credentials", "password"),
	)

	var (
		diags                    diag.Diagnostics
		queryType                = "delete"
		name                     = d.Get("name").(string)
		queryContent             = d.Get("query").(string)
		resourceType             = d.Get("type").(string)
		ignoreAlreadyExists      = d.Get("ignore_already_exists").(bool)
		deleteTopicOnDestroy     = d.Get("delete_topic_on_destroy").(bool)
		terminatePersistentQuery = d.Get("terminate_persistent_query").(bool)
		queryProperties          = d.Get("query_properties").(map[string]interface{})
	)

	eqp := model.NewExecuteQueryRequest(
		name, queryType, queryContent, resourceType,
		ignoreAlreadyExists,
		deleteTopicOnDestroy,
		terminatePersistentQuery,
		cli.MergeWithGlobalProperties(model.NewQueryProperties(queryProperties)),
	)

	err := cli.ExecuteQuery(ctx, eqp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
