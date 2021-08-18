package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelSlack() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelSlackCreate,
		Read:   resourceLaceworkAlertChannelSlackRead,
		Update: resourceLaceworkAlertChannelSlackUpdate,
		Delete: resourceLaceworkAlertChannelSlackDelete,

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
			"slack_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the incoming Slack webhook",
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

func resourceLaceworkAlertChannelSlackCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		slack    = api.NewAlertChannel(d.Get("name").(string),
			api.SlackChannelAlertChannelType,
			api.SlackChannelDataV2{
				SlackUrl: d.Get("slack_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		slack.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.SlackChannelAlertChannelType, slack)
	response, err := lacework.V2.AlertChannels.Create(slack)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.SlackChannelAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.SlackChannelAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.SlackChannelAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelSlackRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.SlackChannelAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetSlackChannel(d.Id())
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
	d.Set("slack_url", response.Data.Data.SlackUrl)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.SlackChannelAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		slack    = api.NewAlertChannel(d.Get("name").(string),
			api.SlackChannelAlertChannelType,
			api.SlackChannelDataV2{
				SlackUrl: d.Get("slack_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		slack.Enabled = 0
	}

	slack.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.SlackChannelAlertChannelType, slack)
	response, err := lacework.V2.AlertChannels.UpdateSlackChannel(slack)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.SlackChannelAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.SlackChannelAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.SlackChannelAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelSlackDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.SlackChannelAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.SlackChannelAlertChannelType, d.Id())
	return nil
}
