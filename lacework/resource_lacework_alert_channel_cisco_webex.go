package lacework

import (
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

func resourceLaceworkAlertChannelCiscoWebexCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		webex    = api.NewAlertChannel(d.Get("name").(string),
			api.CiscoSparkWebhookAlertChannelType,
			api.CiscoSparkWebhookDataV2{
				Webhook: d.Get("webhook_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		webex.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.CiscoSparkWebhookAlertChannelType, webex)
	response, err := lacework.V2.AlertChannels.Create(webex)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.CiscoSparkWebhookAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetCiscoSparkWebhook(d.Id())
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
	d.Set("webhook_url", response.Data.Data.Webhook)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.CiscoSparkWebhookAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		webex    = api.NewAlertChannel(d.Get("name").(string),
			api.CiscoSparkWebhookAlertChannelType,
			api.CiscoSparkWebhookDataV2{
				Webhook: d.Get("webhook_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		webex.Enabled = 0
	}

	webex.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.CiscoSparkWebhookAlertChannelType, webex)
	response, err := lacework.V2.AlertChannels.UpdateCiscoSparkWebhook(webex)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelCiscoWebexDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.CiscoSparkWebhookAlertChannelType, d.Id())
	return nil
}
