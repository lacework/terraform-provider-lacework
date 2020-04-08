package lacework

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwlogger"
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
			"lacework_integration_gcp_cfg":   resourceLaceworkIntegrationGcpCfg(),
			"lacework_integration_gcp_at":    resourceLaceworkIntegrationGcpAt(),
			"lacework_integration_aws_cfg":   resourceLaceworkIntegrationAwsCfg(),
			"lacework_integration_aws_ct":    resourceLaceworkIntegrationAwsCloudTrail(),
			"lacework_integration_azure_cfg": resourceLaceworkIntegrationAzureCfg(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"lacework_api_token": dataSourceLaceworkApiToken(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		err      error
		lacework *api.Client
		account  = d.Get("account").(string)
		key      = d.Get("api_key").(string)
		secret   = d.Get("api_secret").(string)
	)

	// create a new Lacework api client, verify if the terraform command
	// was run with any logging mode, if so, pass it to the lacework client
	logLevel := os.Getenv("TF_LOG")
	if logLevel == "" {
		lacework, err = api.NewClient(account, api.WithApiKeys(key, secret))
	} else {
		// validate that the log level is supported by the api client, if not,
		// use the highest supported level just to help the user troubleshoot
		if !lwlogger.ValidLevel(logLevel) {
			log.Println("[INFO] Unsupported log level for the Lacework provider")
			log.Println("[INFO] Using the 'DEBUG' as the default level")
			logLevel = "DEBUG"
		}

		lacework, err = api.NewClient(
			account,
			api.WithApiKeys(key, secret),
			api.WithLogLevelAndWriter(logLevel, log.Writer()),
		)
	}

	return lacework, err
}
