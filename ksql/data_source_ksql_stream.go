package ksql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStream() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about a KSQL Stream for use in other resources.",
		ReadContext: dataSourceStreamRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the stream.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"topic": {
				Description: "The topic backing the stream.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)

	response, err := client.GetStreamByName(name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("topic", response.Topic)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Name)

	return diags
}
