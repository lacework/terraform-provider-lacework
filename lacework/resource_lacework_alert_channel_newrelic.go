package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelNewRelic() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelNewRelicCreate,
		Read:   resourceLaceworkAlertChannelNewRelicRead,
		Update: resourceLaceworkAlertChannelNewRelicUpdate,
		Delete: resourceLaceworkAlertChannelNewRelicDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAlertChannel,
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
			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The New Relic account ID",
			},
			"insert_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The New Relic Insert API key",
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

func resourceLaceworkAlertChannelNewRelicCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		relic    = api.NewAlertChannel(d.Get("name").(string),
			api.NewRelicInsightsAlertChannelType,
			api.NewRelicInsightsDataV2{
				AccountID: d.Get("account_id").(int),
				InsertKey: d.Get("insert_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		relic.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.NewRelicInsightsAlertChannelType, relic)
	response, err := lacework.V2.AlertChannels.Create(relic)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.NewRelicInsightsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.NewRelicInsightsAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelNewRelicRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetNewRelicInsights(d.Id())
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
	d.Set("account_id", integration.Data.AccountID)
	d.Set("insert_key", integration.Data.InsertKey)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.NewRelicInsightsAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelNewRelicUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		relic    = api.NewAlertChannel(d.Get("name").(string),
			api.NewRelicInsightsAlertChannelType,
			api.NewRelicInsightsDataV2{
				AccountID: d.Get("account_id").(int),
				InsertKey: d.Get("insert_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		relic.Enabled = 0
	}

	relic.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.NewRelicInsightsAlertChannelType, relic)
	response, err := lacework.V2.AlertChannels.UpdateNewRelicInsights(relic)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.NewRelicInsightsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelNewRelicDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.NewRelicInsightsAlertChannelType, d.Id())
	return nil
}
