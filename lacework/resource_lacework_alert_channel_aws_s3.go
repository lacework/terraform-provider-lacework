package lacework

import (
	"log"

	"github.com/lacework/go-sdk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLaceworkAlertChannelAwsS3() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelAwsS3Create,
		Read:   resourceLaceworkAlertChannelAwsS3Read,
		Update: resourceLaceworkAlertChannelAwsS3Update,
		Delete: resourceLaceworkAlertChannelAwsS3Delete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The integration name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"bucket_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the S3 bucket",
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The credentials needed by the integration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ARN of the IAM role",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The external ID of the IAM role",
						},
					},
				},
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modification",
			},
			"intg_guid": {
				Type:     schema.TypeString,
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

func resourceLaceworkAlertChannelAwsS3Create(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		s3       = api.NewAlertChannel(d.Get("name").(string),
			api.AwsS3AlertChannelType,
			api.AwsS3DataV2{
				Credentials: api.AwsS3Credentials{
					ExternalID: d.Get("credentials.0.external_id").(string),
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					BucketArn:  d.Get("bucket_arn").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		s3.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AwsS3AlertChannelType, s3)
	response, err := lacework.V2.AlertChannels.Create(s3)
	if err != nil {
		return err
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.AwsS3AlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			d.SetId("")
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.AwsS3AlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.AwsS3AlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelAwsS3Read(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.AwsS3AlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetAwsS3(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)
	d.Set("bucket_arn", response.Data.Data.Credentials.BucketArn)

	creds := make(map[string]string)
	creds["role_arn"] = response.Data.Data.Credentials.RoleArn
	creds["external_id"] = response.Data.Data.Credentials.ExternalID

	d.Set("credentials", []map[string]string{creds})

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.AwsS3AlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelAwsS3Update(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		s3       = api.NewAlertChannel(d.Get("name").(string),
			api.AwsS3AlertChannelType,
			api.AwsS3DataV2{
				Credentials: api.AwsS3Credentials{
					ExternalID: d.Get("credentials.0.external_id").(string),
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					BucketArn:  d.Get("bucket_arn").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		s3.Enabled = 0
	}

	s3.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsS3AlertChannelType, s3)
	response, err := lacework.V2.AlertChannels.UpdateAwsS3(s3)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.AwsS3AlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.AwsS3AlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.AwsS3AlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelAwsS3Delete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.AwsS3AlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.AwsS3AlertChannelType, d.Id())
	return nil
}
