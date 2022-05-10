package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"webhook_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The url of the external webhook",
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

func resourceLaceworkAlertChannelVictorOpsCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		victor   = api.NewAlertChannel(d.Get("name").(string),
			api.VictorOpsAlertChannelType,
			api.VictorOpsDataV2{
				Url: d.Get("webhook_url").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		victor.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.VictorOpsAlertChannelType, victor)
	response, err := lacework.V2.AlertChannels.Create(victor)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.VictorOpsAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.VictorOpsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.VictorOpsAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelVictorOpsRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.VictorOpsAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetVictorOps(d.Id())
	if err != nil {
		return resourceNotFound(d, err, d.Id())
	}

	integration := response.Data

	if integration.IntgGuid == d.Id() {
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)
		d.Set("webhook_url", integration.Data.Url)

		log.Printf("[INFO] Read %s integration with guid %s\n",
			api.VictorOpsAlertChannelType, integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelVictorOpsUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		victor   = api.NewAlertChannel(d.Get("name").(string),
			api.VictorOpsAlertChannelType,
			api.VictorOpsDataV2{
				Url: d.Get("webhook_url").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		victor.Enabled = 0
	}

	victor.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.VictorOpsAlertChannelType, victor)
	response, err := lacework.V2.AlertChannels.UpdateVictorOps(victor)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.VictorOpsAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.VictorOpsAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.VictorOpsAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelVictorOpsDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.VictorOpsAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.VictorOpsAlertChannelType, d.Id())
	return nil
}
