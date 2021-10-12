package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelSplunk() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelSplunkCreate,
		Read:   resourceLaceworkAlertChannelSplunkRead,
		Update: resourceLaceworkAlertChannelSplunkUpdate,
		Delete: resourceLaceworkAlertChannelSplunkDelete,

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
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hec_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 65536 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 65536 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"event_data": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceLaceworkAlertChannelSplunkCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		splunk   = api.NewSplunkAlertChannel(d.Get("name").(string),
			api.SplunkChannelData{
				Channel:  d.Get("channel").(string),
				HecToken: d.Get("hec_token").(string),
				Host:     d.Get("host").(string),
				Port:     d.Get("port").(int),
				Ssl:      d.Get("ssl").(bool),
				EventData: api.SplunkEventData{
					Index:  d.Get("event_data.0.index").(string),
					Source: d.Get("event_data.0.source").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		splunk.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.SplunkIntegration, splunk)
	response, err := lacework.Integrations.CreateSplunkAlertChannel(splunk)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateSplunkAlertChannelResponse(&response)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.SplunkIntegration, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			d.SetId("")
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.SplunkIntegration, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.SplunkIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelSplunkRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.SplunkIntegration, d.Id())
	response, err := lacework.Integrations.GetSplunkAlertChannel(d.Id())
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
			d.Set("channel", integration.Data.Channel)
			d.Set("hec_token", integration.Data.HecToken)
			d.Set("host", integration.Data.Host)
			d.Set("port", integration.Data.Port)
			d.Set("ssl", integration.Data.Ssl)

			eventData := make(map[string]string)
			eventData["index"] = integration.Data.EventData.Index
			eventData["source"] = integration.Data.EventData.Source

			d.Set("event_data", []map[string]string{eventData})

			log.Printf("[INFO] Read %s integration with guid %s\n",
				api.SplunkIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelSplunkUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		splunk   = api.NewSplunkAlertChannel(d.Get("name").(string),
			api.SplunkChannelData{
				Channel:  d.Get("channel").(string),
				HecToken: d.Get("hec_token").(string),
				Host:     d.Get("host").(string),
				Port:     d.Get("port").(int),
				Ssl:      d.Get("ssl").(bool),
				EventData: api.SplunkEventData{
					Index:  d.Get("event_data.0.index").(string),
					Source: d.Get("event_data.0.source").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		splunk.Enabled = 0
	}

	splunk.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.SplunkIntegration, splunk)
	response, err := lacework.Integrations.UpdateSplunkAlertChannel(splunk)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateSplunkAlertChannelResponse(&response)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.SplunkIntegration, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.SplunkIntegration, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.SplunkIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelSplunkDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.SplunkIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.SplunkIntegration, d.Id())
	return nil
}

func validateSplunkAlertChannelResponse(response *api.SplunkAlertChannelResponse) error {
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
