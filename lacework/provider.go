package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/lacework/go-sdk/api"
)

// Provider returns a Lacework terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_ACCOUNT", nil),
				Description: "Lacework account subdomain of URL (i.e. <ACCOUNT>.lacework.net)",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_API_KEY", nil),
				Description: "Lacework API access key",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_API_SECRET", nil),
				Description: "Lacework API access secret",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"lacework_integration_gcp_cfg": resourceLaceworkIntegrationGcpCfg(),
			"lacework_integration_gcp_at":  resourceLaceworkIntegrationGcpAt(),
			"lacework_integration_aws_cfg": resourceLaceworkIntegrationAwsCfg(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"lacework_api_token": dataSourceLaceworkApiToken(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		account = d.Get("account").(string)
		key     = d.Get("api_key").(string)
		secret  = d.Get("api_secret").(string)
	)

	// create a new Lacework api client
	lacework, err := api.NewClient(
		account,
		api.WithApiKeys(key, secret),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Lacework API client created successfully.")
	return lacework, nil
}
