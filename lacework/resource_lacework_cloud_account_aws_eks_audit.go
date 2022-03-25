package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsEksAudit() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationAwsEksAuditCreate,
		Read:     resourceLaceworkIntegrationAwsEksAuditRead,
		Update:   resourceLaceworkIntegrationAwsEksAuditUpdate,
		Delete:   resourceLaceworkIntegrationAwsEksAuditDelete,
		Schema:   awsEksAuditIntegrationSchema,
		Importer: &schema.ResourceImporter{State: importLaceworkIntegration},
	}
}

var awsEksAuditIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The integration name.",
	},
	"type_name": {
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
	"intg_guid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"sns_arn": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The SNS ARN.",
	},
	"credentials": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Required:    true,
		Description: "The credentials needed by the integration.",
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
	"is_org": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"created_or_updated_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_or_updated_by": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourceLaceworkIntegrationAwsEksAuditCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework        = meta.(*api.Client)
		retries         = d.Get("retries").(int)
		awsEksAuditData = api.AwsEksAuditData{
			SnsArn: d.Get("snsArn").(string),
			Credentials: api.AwsEksAuditCredentials{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
		}
	)

	awsEksAudit := api.NewCloudAccount(d.Get("name").(string),
		api.AwsEksAuditCloudAccount,
		awsEksAuditData,
	)

	if !d.Get("enabled").(bool) {
		awsEksAudit.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AwsEksAuditCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(awsEksAudit)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("error creating %s cloud account integration: %s",
						api.AwsEksAuditCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AwsEksAuditCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"unable to create %s cloud account integration (retrying %d more time(s))",
				api.AwsEksAuditCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		d.SetId(cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)

		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type) // @afiune should we deprecate?
		d.Set("org_level", cloudAccount.IsOrg == 1)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.AwsEksAuditCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsEksAuditRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AwsEksAuditCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsEksAudit(d.Id())
	if err != nil {
		return err
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
		d.Set("snsArn", cloudAccount.Data.SnsArn)

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AwsEksAuditCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsEksAuditUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework        = meta.(*api.Client)
		awsEksAuditData = api.AwsEksAuditData{
			SnsArn: d.Get("snsArn").(string),
			Credentials: api.AwsEksAuditCredentials{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
		}
	)

	awsEksAudit := api.NewCloudAccount(d.Get("name").(string),
		api.AwsEksAuditCloudAccount,
		awsEksAuditData,
	)

	if !d.Get("enabled").(bool) {
		awsEksAudit.Enabled = 0
	}

	awsEksAudit.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsEksAuditCloudAccount.String(), awsEksAudit.IntgGuid)
	response, err := lacework.V2.CloudAccounts.UpdateAwsEksAudit(awsEksAudit)
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

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AwsEksAuditCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsEksAuditDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsEksAuditCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AwsEksAuditCloudAccount.String(), d.Id())
	return nil
}
