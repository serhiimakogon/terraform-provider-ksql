package ksql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStream() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a KSQL Stream resource.",
		CreateContext: resourceStreamCreate,
		ReadContext:   resourceStreamRead,
		DeleteContext: resourceStreamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the stream. Case insensitive.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"query": {
				Description: "The statement to create the stream.",
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	client.url = extractStringValueFromBlock(d, "credentials", "url")
	client.username = extractStringValueFromBlock(d, "credentials", "key")
	client.password = extractStringValueFromBlock(d, "credentials", "secret")

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	query := d.Get("query").(string)

	_, err := client.CreateStream(name, query)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return diags
}

func resourceStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)

	_, err := client.GetStreamByName(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)

	_, err := client.DropStream(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
