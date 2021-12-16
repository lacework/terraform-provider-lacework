package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsGovCloudCT() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsGovCloudCTCreate,
		Read:   resourceLaceworkIntegrationAwsGovCloudCTRead,
		Update: resourceLaceworkIntegrationAwsGovCloudCTUpdate,
		Delete: resourceLaceworkIntegrationAwsGovCloudCTDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
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
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS Account ID",
			},
			"queue_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SQS Queue URL.",
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AWS access key ID",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The AWS secret key for the specified AWS access key",
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

func resourceLaceworkIntegrationAwsGovCloudCTCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsGovCloudCTIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				GovCloudCredentials: &api.AwsGovCloudCreds{
					AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
					SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
					AccountID:       d.Get("account_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AwsGovCloudCTIntegration.String())
		response, err := lacework.Integrations.CreateAws(aws)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.AwsGovCloudCTIntegration.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AwsGovCloudCTIntegration.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.AwsGovCloudCTIntegration.String(), retries,
			))
		}

		log.Printf("[INFO] Verifying server response")
		err = validateAwsIntegrationResponse(&response)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		integration := response.Data[0]
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)

		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.TypeName)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.AwsGovCloudCTIntegration.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsGovCloudCTRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AwsGovCloudCTIntegration.String(), d.Id())
	response, err := lacework.Integrations.GetAws(d.Id())
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)
			d.Set("queue_url", integration.Data.QueueUrl)
			d.Set("account_id", integration.Data.GetAccountID())

			creds := make(map[string]string)
			credentials := integration.Data.GetGovCloudCredentials()
			creds["access_key_id"] = credentials.AccessKeyID
			creds["secret_access_key"] = credentials.SecretAccessKey
			d.Set("credentials", []map[string]string{creds})

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.AwsGovCloudCTIntegration.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsGovCloudCTUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsGovCloudCTIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				GovCloudCredentials: &api.AwsGovCloudCreds{
					AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
					SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
					AccountID:       d.Get("account_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AwsGovCloudCTIntegration.String(), aws)
	response, err := lacework.Integrations.UpdateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
	if err != nil {
		return err
	}

	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.AwsGovCloudCTIntegration.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsGovCloudCTDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AwsGovCloudCTIntegration.String(), d.Id())
	_, err := lacework.Integrations.DeleteAws(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AwsGovCloudCTIntegration.String(), d.Id())
	return nil
}
