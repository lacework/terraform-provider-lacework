package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsAgentlessScanning() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationAwsAgentlessScanningCreate,
		Read:     resourceLaceworkIntegrationAwsAgentlessScanningRead,
		Update:   resourceLaceworkIntegrationAwsAgentlessScanningUpdate,
		Delete:   resourceLaceworkIntegrationAwsAgentlessScanningDelete,
		Schema:   awsAgentlessScanningIntegrationSchema,
		Importer: &schema.ResourceImporter{State: importLaceworkIntegration},
	}
}

var awsAgentlessScanningIntegrationSchema = map[string]*schema.Schema{
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
}

func resourceLaceworkIntegrationAwsAgentlessScanningCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	awsAgentlessScanningData := api.AwsSidekickData{
		ScanFrequency:           d.Get("scan_frequency").(int),
		ScanContainers:          d.Get("scan_containers").(bool),
		ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
	}

	if d.Get("query_text") != nil {
		awsAgentlessScanningData.QueryText = d.Get("query_text").(string)
	}

	awsAgentlessScanning := api.NewCloudAccount(d.Get("name").(string),
		api.AwsSidekickCloudAccount,
		awsAgentlessScanningData,
	)

	if !d.Get("enabled").(bool) {
		awsAgentlessScanning.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AwsSidekickCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAwsSidekick(awsAgentlessScanning)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("Error creating %s cloud account integration: %s",
						api.AwsSidekickCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AwsSidekickCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"Unable to create %s cloud account integration (retrying %d more time(s))",
				api.AwsSidekickCloudAccount.String(), retries,
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
			api.AwsSidekickCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsAgentlessScanningRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AwsSidekickCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsSidekick(d.Id())
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

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AwsSidekickCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsAgentlessScanningUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	awsAgentlessScanningData := api.AwsSidekickData{
		ScanFrequency:           d.Get("scan_frequency").(int),
		ScanContainers:          d.Get("scan_containers").(bool),
		ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
	}

	if d.Get("query_text") != nil {
		awsAgentlessScanningData.QueryText = d.Get("query_text").(string)
	}

	awsAgentlessScanning := api.NewCloudAccount(d.Get("name").(string),
		api.AwsSidekickCloudAccount,
		awsAgentlessScanningData,
	)

	if !d.Get("enabled").(bool) {
		awsAgentlessScanning.Enabled = 0
	}

	awsAgentlessScanning.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsSidekickCloudAccount.String(), awsAgentlessScanning.IntgGuid)
	response, err := lacework.V2.CloudAccounts.UpdateAwsSidekick(awsAgentlessScanning)
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

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AwsSidekickCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsAgentlessScanningDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsSidekickCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AwsSidekickCloudAccount.String(), d.Id())
	return nil
}
