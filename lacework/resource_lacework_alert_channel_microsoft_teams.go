package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelMicrosoftTeams() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelMicrosoftTeamsCreate,
		Read:   resourceLaceworkAlertChannelMicrosoftTeamsRead,
		Update: resourceLaceworkAlertChannelMicrosoftTeamsUpdate,
		Delete: resourceLaceworkAlertChannelMicrosoftTeamsDelete,

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
			"teams_url": {
				Type:     schema.TypeString,
				Required: true,
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
		microsoftTeams = api.NewMicrosoftTeamsAlertChannel(d.Get("name").(string),
			api.MicrosoftTeamsChannelData{
				TeamsURL: d.Get("teams_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		microsoftTeams.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.MicrosoftTeamsChannelIntegration, microsoftTeams)
	response, err := lacework.Integrations.CreateMicrosoftTeamsAlertChannel(microsoftTeams)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateMicrosoftTeamsAlertChannelResponse(&response)
	if err != nil {
		return err
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.MicrosoftTeamsChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.MicrosoftTeamsChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetMicrosoftTeamsAlertChannel(d.Id())
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
			d.Set("teams_url", integration.Data.TeamsURL)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.MicrosoftTeamsChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		microsoftTeams = api.NewMicrosoftTeamsAlertChannel(d.Get("name").(string),
			api.MicrosoftTeamsChannelData{
				TeamsURL: d.Get("teams_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		microsoftTeams.Enabled = 0
	}

	microsoftTeams.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.MicrosoftTeamsChannelIntegration, microsoftTeams)
	response, err := lacework.Integrations.UpdateMicrosoftTeamsAlertChannel(microsoftTeams)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateMicrosoftTeamsAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.MicrosoftTeamsChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelMicrosoftTeamsDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.MicrosoftTeamsChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.MicrosoftTeamsChannelIntegration, d.Id())
	return nil
}

func validateMicrosoftTeamsAlertChannelResponse(response *api.MicrosoftTeamsAlertChannelResponse) error {
	if len(response.Data) == 0 {
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		msg := `
There is more that one integration inside the server response data.

List of integrations:
`
		for _, integration := range response.Data {
			msg = msg + fmt.Sprintf("\t%s: %s\n", integration.IntgGuid, integration.Name)
		}
		msg = msg + unexpectedBehaviorMsg()
		return fmt.Errorf(msg)
	}

	return nil
}
