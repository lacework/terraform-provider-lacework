package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelMicrosoftTeams() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelMicrosoftTeamsCreate,
		Read:   resourceLaceworkAlertChannelMicrosoftTeamsRead,
		Update: resourceLaceworkAlertChannelMicrosoftTeamsUpdate,
		Delete: resourceLaceworkAlertChannelMicrosoftTeamsDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAlertChannel,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert channel",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"webhook_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The webhook url for the integration",
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modification",
			},
			"intg_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the integration",
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

func resourceLaceworkAlertChannelMicrosoftTeamsCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		microsoftTeams = api.NewAlertChannel(d.Get("name").(string),
			api.MicrosoftTeamsAlertChannelType,
			api.MicrosoftTeamsData{
				TeamsURL: d.Get("webhook_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		microsoftTeams.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.MicrosoftTeamsAlertChannelType, microsoftTeams)
	response, err := lacework.V2.AlertChannels.Create(microsoftTeams)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			d.SetId("")
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.MicrosoftTeamsAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetMicrosoftTeams(d.Id())
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
	d.Set("webhook_url", response.Data.Data.TeamsURL)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.MicrosoftTeamsAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		microsoftTeams = api.NewAlertChannel(d.Get("name").(string),
			api.MicrosoftTeamsAlertChannelType,
			api.MicrosoftTeamsData{
				TeamsURL: d.Get("webhook_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		microsoftTeams.Enabled = 0
	}

	microsoftTeams.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.MicrosoftTeamsAlertChannelType, microsoftTeams)
	response, err := lacework.V2.AlertChannels.UpdateMicrosoftTeams(microsoftTeams)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.MicrosoftTeamsAlertChannelType, d.Id())
	return nil
}
