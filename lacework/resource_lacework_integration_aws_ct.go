package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsCloudTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsCloudTrailCreate,
		Read:   resourceLaceworkIntegrationAwsCloudTrailRead,
		Update: resourceLaceworkIntegrationAwsCloudTrailUpdate,
		Delete: resourceLaceworkIntegrationAwsCloudTrailDelete,

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
			"queue_url": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceLaceworkIntegrationAwsCloudTrailCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCloudTrailIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AwsIntegrationCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	// @afiune should we do this if there is sensitive information?
	log.Printf("[INFO] Creating AWS_CFG integration with data:\n%+v\n", aws)
	response, err := lacework.Integrations.CreateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)

	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Created AWS_CFG integration with guid: %v\n", integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading AWS_CFG integration with guid: %v\n", d.Id())
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

			creds := make(map[string]string)
			creds["role_arn"] = integration.Data.Credentials.RoleArn
			creds["external_id"] = integration.Data.Credentials.ExternalID
			d.Set("credentials", []map[string]string{creds})
			d.Set("queue_url", integration.Data.QueueUrl)

			log.Printf("[INFO] Read AWS_CFG integration with guid: %v\n", integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCloudTrailIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AwsIntegrationCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	log.Printf("[INFO] Updating AWS_CFG integration with data:\n%+v\n", aws)
	response, err := lacework.Integrations.UpdateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated AWS_CFG integration with guid: %v\n", d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting AWS_CFG integration with guid: %v\n", d.Id())
	_, err := lacework.Integrations.DeleteAws(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted AWS_CFG integration with guid: %v\n", d.Id())
	return nil
}
