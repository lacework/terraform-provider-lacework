package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkIntegrationAwsCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsCfgCreate,
		Read:   resourceLaceworkIntegrationAwsCfgRead,
		Update: resourceLaceworkIntegrationAwsCfgUpdate,
		Delete: resourceLaceworkIntegrationAwsCfgDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
		},
	}
}

func resourceLaceworkIntegrationAwsCfgCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		aws      = api.NewCloudAccount(d.Get("name").(string),
			api.AwsCfgCloudAccount,
			api.AwsCfgData{
				Credentials: api.AwsCfgCredentials{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			})
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AwsCfgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(aws)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.AwsCfgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AwsCfgCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.AwsCfgCloudAccount.String(), retries,
			))
		}

		integration := response.Data
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)

		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.AwsCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsCfg(d.Id())
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
		credentials := cloudAccount.Data.Credentials
		creds["role_arn"] = credentials.RoleArn
		creds["external_id"] = credentials.ExternalID
		d.Set("credentials", []map[string]string{creds})

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.AwsCfgCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewCloudAccount(d.Get("name").(string),
			api.AwsCfgCloudAccount,
			api.AwsCfgData{
				Credentials: api.AwsCfgCredentials{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			})
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AwsCfgCloudAccount.String(), aws)
	response, err := lacework.V2.CloudAccounts.UpdateAwsCfg(aws)
	if err != nil {
		return err
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	return nil
}
