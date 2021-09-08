package lacework

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwconfig"
	"github.com/lacework/go-sdk/lwlogger"
)

// Provider returns a Lacework schema.Provider
func Provider() *schema.Provider {
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
			"subaccount": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_SUBACCOUNT", nil),
				Description: "The sub-account name inside your organization (org admins only)",
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
			"organization": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LW_ORGANIZATION", nil),
				Description: "Set it to true to access organization level data sets (org admins only)",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"lacework_agent_access_token":            resourceLaceworkAgentAccessToken(),
			"lacework_alert_channel_aws_cloudwatch":  resourceLaceworkAlertChannelAwsCloudWatch(),
			"lacework_alert_channel_aws_s3":          resourceLaceworkAlertChannelAwsS3(),
			"lacework_alert_channel_cisco_webex":     resourceLaceworkAlertChannelCiscoWebex(),
			"lacework_alert_channel_datadog":         resourceLaceworkAlertChannelDatadog(),
			"lacework_alert_channel_email":           resourceLaceworkAlertChannelEmail(),
			"lacework_alert_channel_gcp_pub_sub":     resourceLaceworkAlertChannelGcpPubSub(),
			"lacework_alert_channel_jira_cloud":      resourceLaceworkAlertChannelJiraCloud(),
			"lacework_alert_channel_jira_server":     resourceLaceworkAlertChannelJiraServer(),
			"lacework_alert_channel_newrelic":        resourceLaceworkAlertChannelNewRelic(),
			"lacework_alert_channel_pagerduty":       resourceLaceworkAlertChannelPagerDuty(),
			"lacework_alert_channel_qradar":          resourceLaceworkAlertChannelQRadar(),
			"lacework_alert_channel_microsoft_teams": resourceLaceworkAlertChannelMicrosoftTeams(),
			"lacework_alert_channel_slack":           resourceLaceworkAlertChannelSlack(),
			"lacework_alert_channel_splunk":          resourceLaceworkAlertChannelSplunk(),
			"lacework_alert_channel_service_now":     resourceLaceworkAlertChannelServiceNow(),
			"lacework_alert_channel_victorops":       resourceLaceworkAlertChannelVictorOps(),
			"lacework_alert_channel_webhook":         resourceLaceworkAlertChannelWebhook(),
			"lacework_integration_aws_cfg":           resourceLaceworkIntegrationAwsCfg(),
			"lacework_integration_aws_ct":            resourceLaceworkIntegrationAwsCloudTrail(),
			"lacework_integration_aws_govcloud_cfg":  resourceLaceworkIntegrationAwsGovCloudCfg(),
			"lacework_integration_aws_govcloud_ct":   resourceLaceworkIntegrationAwsGovCloudCT(),
			"lacework_integration_azure_cfg":         resourceLaceworkIntegrationAzureCfg(),
			"lacework_integration_azure_al":          resourceLaceworkIntegrationAzureActivityLog(),
			"lacework_integration_docker_hub":        resourceLaceworkIntegrationDockerHub(),
			"lacework_integration_docker_v2":         resourceLaceworkIntegrationDockerV2(),
			"lacework_integration_ecr":               resourceLaceworkIntegrationEcr(),
			"lacework_integration_gcp_cfg":           resourceLaceworkIntegrationGcpCfg(),
			"lacework_integration_gcp_at":            resourceLaceworkIntegrationGcpAt(),
			"lacework_integration_gar":               resourceLaceworkIntegrationGar(),
			"lacework_integration_gcr":               resourceLaceworkIntegrationGcr(),
			"lacework_integration_ghcr":              resourceLaceworkIntegrationGhcr(),
			"lacework_resource_group_aws":            resourceLaceworkResourceGroupAws(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"lacework_api_token":          dataSourceLaceworkApiToken(),
			"lacework_agent_access_token": dataSourceLaceworkAgentAccessToken(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var (
		diags        diag.Diagnostics
		logLevel     = os.Getenv("TF_LOG")
		profile      = d.Get("profile").(string)
		account      = d.Get("account").(string)
		subaccount   = d.Get("subaccount").(string)
		organization = d.Get("organization").(bool)
		key          = d.Get("api_key").(string)
		secret       = d.Get("api_secret").(string)
		userAgent    = fmt.Sprintf("Terraform/%s", version)
		apiOpts      = []api.Option{api.WithHeader("User-Agent", userAgent)}
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

		if subaccount != "" {
			apiOpts = append(apiOpts, api.WithSubaccount(subaccount))
			apiOpts = append(apiOpts, api.WithApiV2()) // only APIv2 understands sub-accounts
		}

		if organization {
			apiOpts = append(apiOpts, api.WithOrgAccess())
		}

		lw, err := api.NewClient(account, apiOpts...)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Lacework API client",
				Detail:   err.Error(),
			})
		}
		return lw, diags
	}

	if profile == "" {
		profile = "default"
	}

	log.Printf("[INFO] Missing credentials, loading '%s' profile from the Lacework configuration file\n", profile)

	cPath, err := lwconfig.DefaultConfigPath()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// if the Lacework configuration file doesn't exist, we are unable to proceed
	if !fileExist(cPath) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Lacework API client",
			Detail:   providerMisconfiguredErrorMessage(),
		})
		return nil, diags
	}

	profiles, err := lwconfig.LoadProfilesFrom(cPath)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config, ok := profiles[profile]
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Lacework API client",
			Detail: fmt.Sprintf(
				"profile '%s' not found.\n\nTry using the Lacework CLI command 'lacework configure --profile %s'",
				profile, profile),
		})
		return nil, diags
	}

	// Once we have the right credentials loaded from the configuration file,
	// we need to verify if any static setting was provided
	if account == "" {
		account = config.Account
	}

	if subaccount == "" {
		subaccount = config.Subaccount
	}

	if key == "" {
		key = config.ApiKey
	}

	if secret == "" {
		secret = config.ApiSecret
	}

	apiOpts = append(apiOpts, api.WithApiKeys(key, secret))

	if config.Version == 2 {
		// if the config comes back as v2, it means that it is ready to be used
		log.Println("[INFO] Using Lacework APIv2 (loaded from config)")
		apiOpts = append(apiOpts, api.WithApiV2())
	} else {
		// if not, we need to verify that the user is accessing v2 correctly
		log.Println("[INFO] Validating Lacework account information")
		accountVerified, err := verifyPrimaryAccount(account, apiOpts...)
		if err != nil {
			log.Println("[WARN] Unable to validate Lacework account information")
			log.Printf("[WARN] Error: %s\n", err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to validate Lacework account information",
				Detail:   err.Error(),
			})
		}

		log.Println("[INFO] Using Lacework APIv2 (loaded from APIv1 automation)")
		apiOpts = append(apiOpts, api.WithApiV2())

		if accountVerified != account {
			msg := fmt.Sprintf("Overwriting Lacework account to '%s' for APIv2 authentication", accountVerified)
			log.Printf("[WARN] %s\n", msg)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  msg,
				Detail:   providerConfiguredWithV1ConfigMessage(accountVerified),
			})
			account = accountVerified
		}
	}

	if subaccount != "" {
		apiOpts = append(apiOpts, api.WithSubaccount(subaccount))
	}

	if organization {
		apiOpts = append(apiOpts, api.WithOrgAccess())
	}

	lw, err := api.NewClient(account, apiOpts...)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Lacework API client",
			Detail:   err.Error(),
		})
	}
	return lw, diags
}

func verifyPrimaryAccount(account string, opts ...api.Option) (string, error) {
	lwApi, err := api.NewClient(account, opts...)
	if err != nil {
		return account, err
	}

	orgInfo, err := lwApi.Account.GetOrganizationInfo()
	if err != nil {
		return account, err
	}

	if orgInfo.OrgAccount {
		log.Println("[INFO] Organizational account detected")
		return strings.ToLower(orgInfo.AccountName()), nil
	}

	log.Println("[INFO] Account is NOT an organizational account")
	return account, nil
}

func providerMisconfiguredErrorMessage() string {
	return `The Lacework provider is not configured properly. One or more settings are
missing. Refer to the provider documentation for more information:

  https://www.terraform.io/docs/providers/lacework/index.html`
}

func providerConfiguredWithV1ConfigMessage(account string) string {
	return fmt.Sprintf(`
The current account information needs to be changed to the primary account '%s'.

If you are using the Lacework CLI configuration file, verify that you are
running the latest version by running the command:

  lacework version

Refer to the provider documentation for more information:

  https://www.terraform.io/docs/providers/lacework/index.html`, account)
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
