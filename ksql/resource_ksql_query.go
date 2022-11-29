package ksql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-ksql/ksql/client"
)

func resourceQuery() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a KSQL Query resource.",
		CreateContext: resourceQueryCreate,
		ReadContext:   resourceQueryRead,
		DeleteContext: resourceQueryDelete,
		Schema: map[string]*schema.Schema{
			"query": {
				Description: "KSQL query.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"params": {
				Description: "Query parameters to replace. Golang templates format.",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
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
							Required:    true,
							Description: "The KSQL Cluster API URL for your Confluent Cloud cluster.",
							Sensitive:   true,
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The KSQL Cluster API Key for your Confluent Cloud cluster.",
							Sensitive:   true,
						},
						"secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The KSQL Cluster API Secret for your Confluent Cloud cluster.",
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func extractStringValueFromBlock(d *schema.ResourceData, blockName string, attribute string) string {
	// d.Get() will return "" if the key is not present
	return d.Get(fmt.Sprintf("%s.0.%s", blockName, attribute)).(string)
}

func resourceQueryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cli := m.(*client.Client)
	cli.RotateCredentials(
		extractStringValueFromBlock(d, "credentials", "url"),
		extractStringValueFromBlock(d, "credentials", "username"),
		extractStringValueFromBlock(d, "credentials", "password"),
	)

	var (
		diags  diag.Diagnostics
		query  = d.Get("query").(string)
		params = d.Get("params").(map[string]interface{})
	)

	name, err := cli.ExecuteQuery(context.Background(), query, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

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
		diags diag.Diagnostics
		query = d.Get("query").(string)
	)

	tflog.Info(ctx, fmt.Sprintf("Going to read object [%s]", client.ExtractNameFromQuery(query)))

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
		diags diag.Diagnostics
		query = d.Get("query").(string)
	)

	tflog.Info(ctx, fmt.Sprintf("Going to delete object [%s]", client.ExtractNameFromQuery(query)))

	return diags
}
