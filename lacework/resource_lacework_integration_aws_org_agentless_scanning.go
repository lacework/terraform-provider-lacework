package lacework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsOrgAgentlessScanning() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationAwsOrgAgentlessScanningCreate,
		Read:     resourceLaceworkIntegrationAwsOrgAgentlessScanningRead,
		Update:   resourceLaceworkIntegrationAwsOrgAgentlessScanningUpdate,
		Delete:   resourceLaceworkIntegrationAwsOrgAgentlessScanningDelete,
		Schema:   awsOrgAgentlessScanningIntegrationSchema,
		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

var awsOrgAgentlessScanningIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The integration name.",
	},
	"intg_guid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"query_text": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The LQL query text",
	},
	"scan_frequency": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "How often in hours the scan will run in hours.",
	},
	"scan_containers": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether to includes scanning for containers.",
	},
	"scan_host_vulnerabilities": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
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
	"account_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The aws account id",
	},
	"bucket_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The bucket arn",
	},
	"scanning_account": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The scanning aws account id",
	},
	"monitored_accounts": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
		},
		Description: "The list of monitored aws accounts ids or OUs",
	},
	"management_account": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The scanning aws account id",
	},
	"credentials": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Description: "The credentials needed by the integration",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"external_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The external id",
				},
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The role arn",
				},
			},
		},
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
	"created_or_updated_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_or_updated_by": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"org_level": {
		Type:     schema.TypeBool,
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
	"org_account_mappings": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Mapping of AWS accounts to Lacework accounts within a Lacework organization.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_lacework_account": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The default Lacework account name where any non-mapped AWS account will appear",
				},
				"mapping": {
					Type:        schema.TypeSet,
					Required:    true,
					Description: "A map of AWS accounts to Lacework account. This can be specified multiple times to map multiple Lacework accounts.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"lacework_account": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Lacework account name where the Agentless activity from the selected AWS accounts will appear.",
							},
							"aws_accounts": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								MinItems:    1,
								Required:    true,
								Description: "The list of AWS account IDs to map.",
							},
						},
					},
				},
			},
		},
	},
}

func resourceLaceworkIntegrationAwsOrgAgentlessScanningCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	awsOrgAgentlessScanningData := api.AwsSidekickOrgData{
		ScanFrequency:           d.Get("scan_frequency").(int),
		ScanContainers:          d.Get("scan_containers").(bool),
		ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
		AccountID:               d.Get("account_id").(string),
		BucketArn:               d.Get("bucket_arn").(string),
		ScanningAccount:         d.Get("scanning_account").(string),
		MonitoredAccounts:       strings.Join(castAttributeToStringSlice(d, "monitored_accounts"), ", "),
		ManagementAccount:       d.Get("management_account").(string),
		CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
			RoleArn:    d.Get("credentials.0.role_arn").(string),
			ExternalID: d.Get("credentials.0.external_id").(string),
		},
	}

	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d, awsMappingType)
	if !accountMapFile.Empty() {
		accountMapFileBytes, err := json.Marshal(accountMapFile)
		if err != nil {
			return err
		}

		awsOrgAgentlessScanningData.EncodeAccountMappingFile(accountMapFileBytes)
	}

	if d.Get("query_text") != nil {
		awsOrgAgentlessScanningData.QueryText = d.Get("query_text").(string)
	}

	awsOrgAgentlessScanning := api.NewCloudAccount(d.Get("name").(string),
		api.AwsSidekickOrgCloudAccount,
		awsOrgAgentlessScanningData,
	)

	if !d.Get("enabled").(bool) {
		awsOrgAgentlessScanning.Enabled = 0
	}

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AwsSidekickOrgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAwsSidekickOrg(awsOrgAgentlessScanning)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s cloud account integration: %s",
						api.AwsSidekickOrgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AwsSidekickOrgCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s cloud account integration (retrying %d more time(s))",
				api.AwsSidekickOrgCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		d.SetId(cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)

		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type)
		d.Set("org_level", cloudAccount.IsOrg == 1)
		d.Set("server_token", cloudAccount.ServerToken)
		d.Set("uri", cloudAccount.Uri)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.AwsSidekickOrgCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsOrgAgentlessScanningRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AwsSidekickOrgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsSidekickOrg(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	cloudAccount := response.Data
	if cloudAccount.IntgGuid == d.Id() {
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)
		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type)
		d.Set("org_level", cloudAccount.IsOrg == 1)

		creds := make(map[string]string)
		creds["role_arn"] = response.Data.Data.CrossAccountCreds.RoleArn
		creds["external_id"] = response.Data.Data.CrossAccountCreds.ExternalID

		d.Set("credentials", []map[string]string{creds})

		accountMapFileBytes, err := cloudAccount.Data.DecodeAccountMappingFile()
		if err != nil {
			return err
		}

		accountMapFile := new(accountMappingsFile)
		if len(accountMapFileBytes) != 0 {
			// The integration has an account mapping file
			// unmarshal its content into the account mapping struct
			err := json.Unmarshal(accountMapFileBytes, accountMapFile)
			if err != nil {
				return fmt.Errorf("Error decoding organization account mapping: %s", err)
			}

		}

		err = d.Set("org_account_mappings", flattenOrgAccountMappings(accountMapFile, awsMappingType))
		if err != nil {
			return fmt.Errorf("Error flattening organization account mapping: %s", err)
		}

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AwsSidekickOrgCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsOrgAgentlessScanningUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	awsOrgAgentlessScanningData := api.AwsSidekickOrgData{
		ScanFrequency:           d.Get("scan_frequency").(int),
		ScanContainers:          d.Get("scan_containers").(bool),
		ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
		AccountID:               d.Get("account_id").(string),
		BucketArn:               d.Get("bucket_arn").(string),
		ScanningAccount:         d.Get("scanning_account").(string),
		ManagementAccount:       d.Get("management_account").(string),
		MonitoredAccounts:       strings.Join(castAttributeToStringSlice(d, "monitored_accounts"), ", "),
		CrossAccountCreds: api.AwsSidekickCrossAccountCredentials{
			RoleArn:    d.Get("credentials.0.role_arn").(string),
			ExternalID: d.Get("credentials.0.external_id").(string),
		},
	}

	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d, awsMappingType)
	if !accountMapFile.Empty() {
		accountMapFileBytes, err := json.Marshal(accountMapFile)
		if err != nil {
			return err
		}

		awsOrgAgentlessScanningData.EncodeAccountMappingFile(accountMapFileBytes)
	}

	if d.Get("query_text") != nil {
		awsOrgAgentlessScanningData.QueryText = d.Get("query_text").(string)
	}

	awsOrgAgentlessScanning := api.NewCloudAccount(d.Get("name").(string),
		api.AwsSidekickOrgCloudAccount,
		awsOrgAgentlessScanningData,
	)

	if !d.Get("enabled").(bool) {
		awsOrgAgentlessScanning.Enabled = 0
	}

	awsOrgAgentlessScanning.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsSidekickOrgCloudAccount.String(), awsOrgAgentlessScanning.IntgGuid)
	response, err := lacework.V2.CloudAccounts.UpdateAwsSidekickOrg(awsOrgAgentlessScanning)
	if err != nil {
		return err
	}

	cloudAccount := response.Data
	d.Set("name", cloudAccount.Name)
	d.Set("intg_guid", cloudAccount.IntgGuid)
	d.Set("enabled", cloudAccount.Enabled == 1)
	d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
	d.Set("type_name", cloudAccount.Type)
	d.Set("org_level", cloudAccount.IsOrg == 1)

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AwsSidekickOrgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsOrgAgentlessScanningDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsSidekickOrgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AwsSidekickOrgCloudAccount.String(), d.Id())
	return nil
}
