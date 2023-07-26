package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelPagerDuty() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelPagerDutyCreate,
		Read:   resourceLaceworkAlertChannelPagerDutyRead,
		Update: resourceLaceworkAlertChannelPagerDutyUpdate,
		Delete: resourceLaceworkAlertChannelPagerDutyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"integration_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The PagerDuty service integration key",
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation or modification",
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

func resourceLaceworkAlertChannelPagerDutyCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAlertChannel(d.Get("name").(string),
			api.PagerDutyApiAlertChannelType,
			api.PagerDutyApiDataV2{
				IntegrationKey: d.Get("integration_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.PagerDutyApiAlertChannelType, alert)
	response, err := lacework.V2.AlertChannels.Create(alert)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.PagerDutyApiAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.PagerDutyApiAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelPagerDutyRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetPagerDutyApi(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.PagerDutyApiAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelPagerDutyUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAlertChannel(d.Get("name").(string),
			api.PagerDutyApiAlertChannelType,
			api.PagerDutyApiDataV2{
				IntegrationKey: d.Get("integration_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	alert.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.PagerDutyApiAlertChannelType, alert)
	response, err := lacework.V2.AlertChannels.UpdatePagerDutyApi(alert)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.PagerDutyApiAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelPagerDutyDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.PagerDutyApiAlertChannelType, d.Id())
	return nil
}
