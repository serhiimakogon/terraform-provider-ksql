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
		ReadContext: dataSourceStreamsRead,
		Schema: map[string]*schema.Schema{
			"tag": {
				Description: "The tag to filter the streams. Case sensitive.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"topic": {
				Description: "The topic to filter the streams. Case sensitive.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
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

	client.url = extractStringValueFromBlock(d, "credentials", "url")
	client.apiKey = extractStringValueFromBlock(d, "credentials", "key")
	client.apiSecret = extractStringValueFromBlock(d, "credentials", "secret")

	var diags diag.Diagnostics
	var streamsData []Stream
	var err error

	tag := d.Get("tag").(string)
	topic := d.Get("topic").(string)

	if tag != "" && topic != "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "only one filter type is allowed.",
		})
		return diags

	} else if tag != "" {
		streamsData, err = client.GetStreamsByTag(tag)
		if err != nil {
			return diag.FromErr(err)
		}

	} else if topic != "" {
		streamsData, err = client.GetStreamsByTopic(topic)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		streamsData, err = client.ListStreams()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	streams := flattenStreamsData(streamsData)

	err = d.Set("streams", streams)
	if err != nil {
		return diag.FromErr(err)
	}

	id := uuid.New()
	d.SetId(id.String())

	return diags
}

func flattenStreamsData(streamsData []Stream) []map[string]interface{} {
	streams := make([]map[string]interface{}, len(streamsData))

	for i, s := range streamsData {
		stream := make(map[string]interface{})

		stream["name"] = s.Name
		stream["topic"] = s.Topic

		streams[i] = stream
	}
	return streams
}
