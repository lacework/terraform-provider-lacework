package lacework

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwlogger"
)

// Provider returns a Lacework terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_PROFILE", nil),
				Description: "Lacework profile name to use, profiles are configured at $HOME/.lacework.toml via the Lacework CLI",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_ACCOUNT", nil),
				Description: "Lacework account subdomain of URL (i.e. <ACCOUNT>.lacework.net)",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_API_KEY", nil),
				Description: "Lacework API access key",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_API_SECRET", nil),
				Description: "Lacework API access secret",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"lacework_alert_channel_aws_cloudwatch": resourceLaceworkAlertChannelAwsCloudWatch(),
			"lacework_alert_channel_jira_cloud":     resourceLaceworkAlertChannelJiraCloud(),
			"lacework_alert_channel_jira_server":    resourceLaceworkAlertChannelJiraServer(),
			"lacework_alert_channel_pagerduty":      resourceLaceworkAlertChannelPagerDuty(),
			"lacework_alert_channel_slack":          resourceLaceworkAlertChannelSlack(),
			"lacework_integration_aws_cfg":          resourceLaceworkIntegrationAwsCfg(),
			"lacework_integration_aws_ct":           resourceLaceworkIntegrationAwsCloudTrail(),
			"lacework_integration_azure_cfg":        resourceLaceworkIntegrationAzureCfg(),
			"lacework_integration_azure_al":         resourceLaceworkIntegrationAzureActivityLog(),
			"lacework_integration_docker_hub":       resourceLaceworkIntegrationDockerHub(),
			"lacework_integration_docker_v2":        resourceLaceworkIntegrationDockerV2(),
			"lacework_integration_ecr":              resourceLaceworkIntegrationEcr(),
			"lacework_integration_gcp_cfg":          resourceLaceworkIntegrationGcpCfg(),
			"lacework_integration_gcp_at":           resourceLaceworkIntegrationGcpAt(),
			"lacework_integration_gcr":              resourceLaceworkIntegrationGcr(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"lacework_api_token": dataSourceLaceworkApiToken(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		err       error
		logLevel  = os.Getenv("TF_LOG")
		profile   = d.Get("profile").(string)
		account   = d.Get("account").(string)
		key       = d.Get("api_key").(string)
		secret    = d.Get("api_secret").(string)
		userAgent = fmt.Sprintf("Terraform/%s", version)
		apiOpts   = []api.Option{api.WithHeader("User-Agent", userAgent)}
	)

	// validate that the log level is supported by the api client, if not,
	// use the highest supported level just to help the user troubleshoot
	if logLevel != "" {
		if !lwlogger.ValidLevel(logLevel) {
			log.Println("[INFO] Unsupported log level for the Lacework provider")
			log.Println("[INFO] Using the 'DEBUG' as the default level")
			logLevel = "DEBUG"
		}
		apiOpts = append(apiOpts, api.WithLogLevelAndWriter(logLevel, log.Writer()))
	}

	if account != "" && key != "" && secret != "" {
		apiOpts = append(apiOpts, api.WithApiKeys(key, secret))
		return api.NewClient(account, apiOpts...)
	}

	if profile == "" {
		profile = "default"
	}

	log.Printf("[INFO] Missing credentials, loading '%s' profile from the Lacework configuration file\n", profile)

	// read config file $HOME/.lacework.toml
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var (
		profiles = Profiles{}
		cPath    = path.Join(home, ".lacework.toml")
	)

	// if the Lacework configuration file doesn't exist, we are unable to proceed
	if !fileExist(cPath) {
		return nil, errors.New(providerMisconfiguredErrorMessage())
	}

	if _, err := toml.DecodeFile(cPath, &profiles); err != nil {
		return nil, errors.Wrap(err, "unable to decode profiles from config")
	}

	creds, ok := profiles[profile]
	if !ok {
		return nil, errors.Errorf(
			"profile '%s' not found.\n\nTry using the Lacework CLI command 'lacework configure --profile %s'.",
			profile, profile)
	}

	// Once we have the right credentials loaded from the configuration file,
	// we need to verify if any static setting was provided
	if account == "" {
		account = creds.Account
	}
	if key == "" {
		key = creds.ApiKey
	}
	if secret == "" {
		secret = creds.ApiSecret
	}

	apiOpts = append(apiOpts, api.WithApiKeys(key, secret))
	return api.NewClient(account, apiOpts...)
}

func providerMisconfiguredErrorMessage() string {
	return `

The Lacework provider is not configured properly. One or more settings are
missing. Refer to the provider documentation for more information:

  https://www.terraform.io/docs/providers/lacework/index.html`
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type Profiles map[string]credsDetails

type credsDetails struct {
	Account   string `toml:"account"`
	ApiKey    string `toml:"api_key"`
	ApiSecret string `toml:"api_secret"`
}
