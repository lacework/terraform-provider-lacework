package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelVictorOps() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelVictorOpsCreate,
		Read:   resourceLaceworkAlertChannelVictorOpsRead,
		Update: resourceLaceworkAlertChannelVictorOpsUpdate,
		Delete: resourceLaceworkAlertChannelVictorOpsDelete,

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

func resourceLaceworkAlertChannelVictorOpsCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		victor   = api.NewVictorOpsAlertChannel(d.Get("name").(string),
			api.VictorOpsChannelData{
				WebhookURL: d.Get("webhook_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		victor.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.VictorOpsChannelIntegration, victor)
	response, err := lacework.Integrations.CreateVictorOpsAlertChannel(victor)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateVictorOpsAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.VictorOpsChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelVictorOpsRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.VictorOpsChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetVictorOpsAlertChannel(d.Id())
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
				api.VictorOpsChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelVictorOpsUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		victor   = api.NewVictorOpsAlertChannel(d.Get("name").(string),
			api.VictorOpsChannelData{
				WebhookURL: d.Get("webhook_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		victor.Enabled = 0
	}

	victor.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.VictorOpsChannelIntegration, victor)
	response, err := lacework.Integrations.UpdateVictorOpsAlertChannel(victor)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateVictorOpsAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.VictorOpsChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelVictorOpsDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.VictorOpsChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.VictorOpsChannelIntegration, d.Id())
	return nil
}

func validateVictorOpsAlertChannelResponse(response *api.VictorOpsAlertChannelResponse) error {
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