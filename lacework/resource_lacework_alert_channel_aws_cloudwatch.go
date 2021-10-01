package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelAwsCloudWatch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelAwsCloudWatchCreate,
		Read:   resourceLaceworkAlertChannelAwsCloudWatchRead,
		Update: resourceLaceworkAlertChannelAwsCloudWatchUpdate,
		Delete: resourceLaceworkAlertChannelAwsCloudWatchDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The integration name",
			},
			"intg_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The integration unique identifier",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"event_bus_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of your AWS CloudWatch event bus",
			},
			"group_issues_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Events",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Events", "Resources":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'Events' or 'Resources' (default: Events)", key,
							),
						}
					}
				},
				Description: "Defines how Lacework compliance events get grouped. Must be one of Events or Resources. Defaults to Events",
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modification",
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

func resourceLaceworkAlertChannelAwsCloudWatchCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAlertChannel(d.Get("name").(string),
			api.CloudwatchEbAlertChannelType,
			api.CloudwatchEbDataV2{
				EventBusArn:   d.Get("event_bus_arn").(string),
				IssueGrouping: d.Get("group_issues_by").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AwsCloudWatchIntegration, alert)
	response, err := lacework.V2.AlertChannels.Create(alert)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.CloudwatchEbAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.CloudwatchEbAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.CloudwatchEbAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.AwsCloudWatchIntegration, d.Id())
	response, err := lacework.V2.AlertChannels.GetCloudwatchEb(d.Id())
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

	d.Set("event_bus_arn", response.Data.Data.EventBusArn)
	d.Set("group_issues_by", response.Data.Data.IssueGrouping)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.CloudwatchEbAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAlertChannel(d.Get("name").(string),
			api.CloudwatchEbAlertChannelType,
			api.AwsCloudWatchData{
				EventBusArn:   d.Get("event_bus_arn").(string),
				IssueGrouping: d.Get("group_issues_by").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	alert.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsCloudWatchIntegration, alert)
	response, err := lacework.V2.AlertChannels.UpdateCloudwatchEb(alert)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.CloudwatchEbAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.CloudwatchEbAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.CloudwatchEbAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.CloudwatchEbAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.CloudwatchEbAlertChannelType, d.Id())
	return nil
}
