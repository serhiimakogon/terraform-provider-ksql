package ksql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-ksql/ksql/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
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
		ResourcesMap: map[string]*schema.Resource{
			"ksql_query": resourceQuery(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Initializing Terraform Provider for KSQL")

	var (
		url                 = d.Get("url").(string)
		username            = d.Get("username").(string)
		password            = d.Get("password").(string)
		autoOffsetResetMode = d.Get("auto_offset_reset").(string)
	)

	if autoOffsetResetMode != "earliest" && autoOffsetResetMode != "latest" {
		return nil, diag.Errorf("invalid auto_offset_reset mode: %s", autoOffsetResetMode)
	}

	return client.New(url, username, password, autoOffsetResetMode), nil
}
