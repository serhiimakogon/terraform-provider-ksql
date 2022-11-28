package ksql

import (
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
	}
}
