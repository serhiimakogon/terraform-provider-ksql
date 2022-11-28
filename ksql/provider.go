package ksql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ksql_stream": resourceStream(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ksql_stream":  dataSourceStream(),
			"ksql_streams": dataSourceStreams(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	return NewEmptyClient(ctx), diags
}
