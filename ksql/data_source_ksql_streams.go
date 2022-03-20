package ksql

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStreams() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about KSQL Streams for use in other resources.",
		ReadContext: dataSourceStreamRead,
		Schema: map[string]*schema.Schema{
			"tag": {
				Description: "The tag to filter the streams.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"streams": {
				Description: "The streams found.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the stream.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"topic": {
							Description: "The topic backing the stream.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceStreamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	tag := d.Get("tag").(string)

	response, err := client.GetStreamsByTag(tag)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("streams", response)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	return diags
}
