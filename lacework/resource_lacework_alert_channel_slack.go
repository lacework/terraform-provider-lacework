package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
			"slack_url": {
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

func resourceLaceworkAlertChannelSlackCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		slack    = api.NewSlackAlertChannel(d.Get("name").(string),
			api.SlackChannelData{
				SlackUrl: d.Get("slack_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		slack.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.SlackChannelIntegration, slack)
	response, err := lacework.Integrations.CreateSlackAlertChannel(slack)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateSlackAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.SlackChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelSlackRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.SlackChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetSlackAlertChannel(d.Id())
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
			d.Set("slack_url", integration.Data.SlackUrl)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.SlackChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		slack    = api.NewSlackAlertChannel(d.Get("name").(string),
			api.SlackChannelData{
				SlackUrl: d.Get("slack_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		slack.Enabled = 0
	}

	slack.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.SlackChannelIntegration, slack)
	response, err := lacework.Integrations.UpdateSlackAlertChannel(slack)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateSlackAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.SlackChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelSlackDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.SlackChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.SlackChannelIntegration, d.Id())
	return nil
}

// validateSlackAlertChannelResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateSlackAlertChannelResponse(response *api.SlackAlertChannelResponse) error {
	if len(response.Data) == 0 {
		// @afiune this edge case should never happen, if we land here it means that
		// something went wrong in the server side of things (Lacework API), so let
		// us inform that to our users
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		// @afiune if we are creating a single integration and the server returns
		// more than one integration inside the 'data' field, it is definitely another
		// edge case that should never happen
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
