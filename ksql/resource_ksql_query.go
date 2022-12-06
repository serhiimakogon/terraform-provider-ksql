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
							DefaultFunc: schema.EnvDefaultFunc("KSQL_URL", ""),
							Description: "The KSQL URL.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("KSQL_USERNAME", ""),
							Description: "The KSQL username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("KSQL_PASSWORD", ""),
							Description: "The KSQL password.",
						},
					},
				},
			},
		},
	}
}

func extractStringValueFromBlock(d *schema.ResourceData, blockName string, attribute string) string {
	// d.Get() will return "" if the key is not present
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
		diags diag.Diagnostics
		name  = d.Get("name").(string)
		qType = d.Get("type").(string)
		query = d.Get("query").(string)
	)

	id, err := cli.ExecuteQuery(context.Background(), name, qType, query)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

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
