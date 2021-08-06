package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelCiscoWebex() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelCiscoWebexCreate,
		Read:   resourceLaceworkAlertChannelCiscoWebexRead,
		Update: resourceLaceworkAlertChannelCiscoWebexUpdate,
		Delete: resourceLaceworkAlertChannelCiscoWebexDelete,

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
			"webhook_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation",
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

func resourceLaceworkAlertChannelCiscoWebexCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		webex    = api.NewCiscoWebexAlertChannel(d.Get("name").(string),
			api.CiscoWebexChannelData{
				WebhookURL: d.Get("webhook_url").(string),
			},
		)
		testIntegration = d.Get("test_integration").(bool)
	)
	if !d.Get("enabled").(bool) {
		webex.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.CiscoWebexChannelIntegration, webex)
	response, err := lacework.Integrations.CreateCiscoWebexAlertChannel(webex)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateCiscoWebexAlertChannelResponse(&response)
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

	if testIntegration {
		log.Printf("[INFO] Testing %s integration for guid:%s\n", api.DatadogChannelIntegration, d.Id())
		err := VerifyAlertChannel(d.Id(), lacework)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid: %s successfully \n", api.DatadogChannelIntegration, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.CiscoWebexChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.CiscoWebexChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetCiscoWebexAlertChannel(d.Id())
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
			d.Set("webhook_url", integration.Data.WebhookURL)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.CiscoWebexChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		webex    = api.NewCiscoWebexAlertChannel(d.Get("name").(string),
			api.CiscoWebexChannelData{
				WebhookURL: d.Get("webhook_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		webex.Enabled = 0
	}

	webex.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.CiscoWebexChannelIntegration, webex)
	response, err := lacework.Integrations.UpdateCiscoWebexAlertChannel(webex)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateCiscoWebexAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.CiscoWebexChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.CiscoWebexChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.CiscoWebexChannelIntegration, d.Id())
	return nil
}

func validateCiscoWebexAlertChannelResponse(response *api.CiscoWebexAlertChannelResponse) error {
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
