package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelPagerDuty() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelPagerDutyCreate,
		Read:   resourceLaceworkAlertChannelPagerDutyRead,
		Update: resourceLaceworkAlertChannelPagerDutyUpdate,
		Delete: resourceLaceworkAlertChannelPagerDutyDelete,

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
			"integration_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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

func resourceLaceworkAlertChannelPagerDutyCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewPagerDutyAlertChannel(d.Get("name").(string),
			api.PagerDutyData{
				IntegrationKey: d.Get("integration_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.PagerDutyIntegration, alert)
	response, err := lacework.Integrations.CreatePagerDutyAlertChannel(alert)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validatePagerDutyAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.PagerDutyIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelPagerDutyRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.PagerDutyIntegration, d.Id())
	response, err := lacework.Integrations.GetPagerDutyAlertChannel(d.Id())
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
			d.Set("integration_key", integration.Data.IntegrationKey)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.PagerDutyIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelPagerDutyUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewPagerDutyAlertChannel(d.Get("name").(string),
			api.PagerDutyData{
				IntegrationKey: d.Get("integration_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	alert.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.PagerDutyIntegration, alert)
	response, err := lacework.Integrations.UpdatePagerDutyAlertChannel(alert)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validatePagerDutyAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.PagerDutyIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelPagerDutyDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.PagerDutyIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.PagerDutyIntegration, d.Id())
	return nil
}

// validatePagerDutyAlertChannelResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validatePagerDutyAlertChannelResponse(response *api.PagerDutyAlertChannelResponse) error {
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
