package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAzureAgentlessScanning() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAzureAgentlessScanningCreate,
		Read:   resourceLaceworkIntegrationAzureAgentlessScanningRead,
		Update: resourceLaceworkIntegrationAzureAgentlessScanningUpdate,
		Delete: resourceLaceworkIntegrationAzureAgentlessScanningDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The integration name.",
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration.",
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of attempts to create the external integration.",
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client Id from credentials file.",
						},
						"client_secret": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// we don't compare this field for security reasons. so here only compare other fields
								return !d.HasChanges(
									"name",
									"integration_level",
									"tenant_id",
									"scanning_subscription_id",
									"scanning_resource_group_name",
									"storage_account_url",
									"enabled",
									"credentials.0.client_id",
								)
							},
							Description: "Client secret from credentials file.",
						},
					},
				},
			},
			"integration_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.AzureSubscriptionIntegration,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch strings.ToUpper(value.(string)) {
					case api.AzureSubscriptionIntegration,
						api.AzureTenantIntegration:
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf("%s: can only be either '%s' or '%s'",
								key,
								api.AzureSubscriptionIntegration,
								api.AzureTenantIntegration),
						}
					}
				},
				Description: "Integration level - TENANT / SUBSCRIPTION.",
			},
			"scanning_subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the subscription where LW scanner is deployed.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tenant where LW scanner is deployed.",
			},
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_or_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blob_container_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "blob container containing analysis results shared with Lacework platform.",
			},
			"scanning_resource_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the resource group where the scanner runs.",
			},
			"storage_account_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the storage account, in the format of 'https://<account>.blob.core.windows.net'",
			},
			"scan_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     24,
				Description: "How often in hours the scan will run in hours.",
			},
			"scan_containers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to includes scanning for containers.",
			},
			"scan_host_vulnerabilities": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to includes scanning for host vulnerabilities.",
			},
			"scan_multi_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to scan secondary volumes (true) or only root volumes (false)",
			},
			"scan_stopped_instances": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to scan stopped instances (true)",
			},
			"query_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The LQL query text.",
			},
			"subscriptions_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Default:     nil,
				Description: "List of subscriptions to specifically include/exclude.",
			},
		},
	}
}

func resourceLaceworkIntegrationAzureAgentlessScanningCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework         = meta.(*api.Client)
		retries          = d.Get("retries").(int)
		integrationLevel = api.AzureSubscriptionIntegration
	)

	if strings.ToUpper(
		d.Get("integration_level").(string),
	) == api.AzureTenantIntegration {
		integrationLevel = api.AzureTenantIntegration
	}
	log.Printf("[INFO] Creating %s integration\n", api.AzureSidekickCloudAccount.String())

	data := api.NewCloudAccount(d.Get("name").(string),
		api.AzureSidekickCloudAccount,
		api.AzureSidekickData{
			IntegrationLevel: integrationLevel,
			Credentials: api.AzureSidekickCredentials{
				ClientId:     d.Get("credentials.0.client_id").(string),
				ClientSecret: d.Get("credentials.0.client_secret").(string),
			},
			BlobContainerName:         d.Get("blob_container_name").(string),
			ScanningResourceGroupName: d.Get("scanning_resource_group_name").(string),
			StorageAccountUrl:         d.Get("storage_account_url").(string),
			ScanningSubscriptionId:    d.Get("scanning_subscription_id").(string),
			TenantId:                  d.Get("tenant_id").(string),
			ScanFrequency:             d.Get("scan_frequency").(int),
			ScanContainers:            d.Get("scan_containers").(bool),
			ScanHostVulnerabilities:   d.Get("scan_host_vulnerabilities").(bool),
			ScanMultiVolume:           d.Get("scan_multi_volume").(bool),
			ScanStoppedInstances:      d.Get("scan_stopped_instances").(bool),
			QueryText:                 d.Get("query_text").(string),
			SubscriptionsList:         strings.Join(castAttributeToStringSlice(d, "subscriptions_list"), ", "),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AzureSidekickCloudAccount.String())
		log.Printf("[INFO] Creating %v integration\n", data)
		response, err := lacework.V2.CloudAccounts.CreateAzureSidekick(data)

		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.AzureSidekickCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AzureSidekickCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.AzureSidekickCloudAccount.String(), retries,
			))
		}

		integration := response.Data
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("server_token", integration.ServerToken)
		d.Set("uri", integration.Uri)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.AzureSidekickCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAzureAgentlessScanningRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AzureSidekickCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAzureSidekick(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	integration := response.Data
	if integration.IntgGuid == d.Id() {
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("integration_level", integration.Type)

		creds := make(map[string]string)
		creds["client_id"] = integration.Data.Credentials.ClientId
		d.Set("credentials", []map[string]string{creds})

		d.Set("integration_level", integration.Data.IntegrationLevel)
		d.Set("scanning_subscription_id", integration.Data.ScanningSubscriptionId)
		d.Set("tenant_id", integration.Data.TenantId)
		d.Set("blob_container_name", integration.Data.BlobContainerName)
		d.Set("storage_account_url", integration.Data.StorageAccountUrl)
		d.Set("scanning_resource_group_name", integration.Data.ScanningResourceGroupName)
		d.Set("scan_frequency", integration.Data.ScanFrequency)
		d.Set("scan_containers", integration.Data.ScanContainers)
		d.Set("scan_host_vulnerabilities", integration.Data.ScanHostVulnerabilities)
		d.Set("query_text", integration.Data.QueryText)
		d.Set("uri", integration.Uri)

		subscriptions_list := strings.Split(integration.Data.SubscriptionsList, ",")
		if integration.Data.SubscriptionsList != "" && len(subscriptions_list) > 0 {
			var trimmed_subscriptions_list []string
			for _, elem := range subscriptions_list {
				trimmed_subscriptions_list = append(trimmed_subscriptions_list, strings.TrimSpace(elem))
			}
			d.Set("subscriptions_list", trimmed_subscriptions_list)
		}

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.AzureSidekickCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationAzureAgentlessScanningUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework         = meta.(*api.Client)
		integrationLevel = api.AzureSubscriptionIntegration
	)

	if strings.ToUpper(d.Get("integration_level").(string)) == api.AzureTenantIntegration {
		integrationLevel = api.AzureTenantIntegration
	}

	data := api.NewCloudAccount(d.Get("name").(string),
		api.AzureSidekickCloudAccount,
		api.AzureSidekickData{
			ScanningSubscriptionId: d.Get("scanning_subscription_id").(string),
			TenantId:               d.Get("tenant_id").(string),
			IntegrationLevel:       integrationLevel,
			Credentials: api.AzureSidekickCredentials{
				ClientId:     d.Get("credentials.0.client_id").(string),
				ClientSecret: d.Get("credentials.0.client_secret").(string),
			},
			BlobContainerName:         d.Get("blob_container_name").(string),
			ScanningResourceGroupName: d.Get("scanning_resource_group_name").(string),
			StorageAccountUrl:         d.Get("storage_account_url").(string),
			ScanFrequency:             d.Get("scan_frequency").(int),
			ScanContainers:            d.Get("scan_containers").(bool),
			ScanHostVulnerabilities:   d.Get("scan_host_vulnerabilities").(bool),
			ScanMultiVolume:           d.Get("scan_multi_volume").(bool),
			ScanStoppedInstances:      d.Get("scan_stopped_instances").(bool),
			QueryText:                 d.Get("query_text").(string),
			SubscriptionsList:         strings.Join(castAttributeToStringSlice(d, "subscriptions_list"), ", "),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AzureSidekickCloudAccount.String(), data)
	response, err := lacework.V2.CloudAccounts.UpdateAzureSidekick(data)
	if err != nil {
		return err
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("integration_level", integration.Data.IntegrationLevel)
	d.Set("tenant_id", integration.Data.TenantId)
	d.Set("scanning_subscription_id", integration.Data.ScanningSubscriptionId)
	d.Set("blob_container_name", integration.Data.BlobContainerName)
	d.Set("storage_account_url", integration.Data.StorageAccountUrl)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("server_token", integration.ServerToken)
	d.Set("uri", integration.Uri)
	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.AzureSidekickCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAzureAgentlessScanningDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AzureSidekickCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AzureSidekickCloudAccount.String(), d.Id())
	return nil
}

// TODO: generate the documentation
