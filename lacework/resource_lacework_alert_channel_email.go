package lacework

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelEmailCreate,
		Read:   resourceLaceworkAlertChannelEmailRead,
		Update: resourceLaceworkAlertChannelEmailUpdate,
		Delete: resourceLaceworkAlertChannelEmailDelete,

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
			"recipients": {
				Type:        schema.TypeList,
				MinItems:    1,
				Required:    true,
				Description: "List of email addresses that you want to receive the alerts",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
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

func resourceLaceworkAlertChannelEmailCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		emailAlertChan = api.NewAlertChannel(d.Get("name").(string),
			api.EmailUserAlertChannelType,
			api.EmailUserData{
				ChannelProps: api.EmailUserChannelProps{
					Recipients: castAttributeToStringSlice(d, "recipients"),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		emailAlertChan.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.EmailUserAlertChannelType, emailAlertChan)
	response, err := lacework.V2.AlertChannels.Create(emailAlertChan)
	if err != nil {
		return err
	}

	d.SetId(response.Data.IntgGuid)
	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.EmailUserAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.EmailUserAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid: %s\n", api.EmailUserAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelEmailRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %s\n", api.EmailUserAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetEmailUser(d.Id())
	if err != nil {
		return resourceNotFound(d, err, d.Id())
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)
	// @afiune TODO
	d.Set("recipients", response.Data.Data.ChannelProps.Recipients)

	log.Printf("[INFO] Read %s integration with guid: %s\n",
		api.EmailUserAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		emailAlertChan = api.NewAlertChannel(d.Get("name").(string),
			api.EmailUserAlertChannelType,
			api.EmailUserData{
				ChannelProps: api.EmailUserChannelProps{
					Recipients: castAttributeToStringSlice(d, "recipients"),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		emailAlertChan.Enabled = 0
	}

	emailAlertChan.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.EmailUserAlertChannelType, emailAlertChan)
	response, err := lacework.V2.AlertChannels.UpdateEmailUser(emailAlertChan)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.EmailUserAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid: %s successfully\n", api.EmailUserAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.EmailUserAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelEmailDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %s\n", api.EmailUserAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %s\n", api.EmailUserAlertChannelType, d.Id())
	return nil
}
