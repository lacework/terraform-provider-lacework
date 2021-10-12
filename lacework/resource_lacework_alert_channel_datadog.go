package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelDatadog() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelDatadogCreate,
		Read:   resourceLaceworkAlertChannelDatadogRead,
		Update: resourceLaceworkAlertChannelDatadogUpdate,
		Delete: resourceLaceworkAlertChannelDatadogDelete,

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
			"datadog_site": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(api.DatadogSiteCom),
				Description: "Where to store your logs, either the US or Europe",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case string(api.DatadogSiteEu), string(api.DatadogSiteCom):
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either '%s' or '%s'", key,
								api.DatadogSiteEu, api.DatadogSiteCom,
							),
						}
					}
				},
			},
			"datadog_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(api.DatadogServiceLogsDetails),
				Description: "The level of detail of logs or event stream",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case string(api.DatadogServiceLogsDetails), string(api.DatadogServiceLogsSummary), string(api.DatadogServiceEventsSummary):
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either '%s', '%s' or '%s'", key,
								api.DatadogServiceLogsDetails, api.DatadogServiceLogsSummary, api.DatadogServiceEventsSummary,
							),
						}
					}
				},
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The Datadog api key required to submit metrics and events to Datadog",
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modifications",
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

func resourceLaceworkAlertChannelDatadogCreate(d *schema.ResourceData, meta interface{}) error {
	site, _ := api.DatadogSite(d.Get("datadog_site").(string))
	service, _ := api.DatadogService(d.Get("datadog_service").(string))

	var (
		lacework = meta.(*api.Client)
		datadog  = api.NewAlertChannel(d.Get("name").(string),
			api.DatadogAlertChannelType,
			api.DatadogDataV2{
				DatadogSite: site,
				DatadogType: service,
				ApiKey:      d.Get("api_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		datadog.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.DatadogAlertChannelType, datadog)
	response, err := lacework.V2.AlertChannels.Create(datadog)

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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.DatadogAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.DatadogAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid: %s\n", api.DatadogAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelDatadogRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.DatadogAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetDatadog(d.Id())
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
	d.Set("datadog_site", response.Data.Data.DatadogSite)
	d.Set("datadog_service", response.Data.Data.DatadogType)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.DatadogAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelDatadogUpdate(d *schema.ResourceData, meta interface{}) error {
	site, _ := api.DatadogSite(d.Get("datadog_site").(string))
	service, _ := api.DatadogService(d.Get("datadog_service").(string))

	var (
		lacework = meta.(*api.Client)
		datadog  = api.NewAlertChannel(d.Get("name").(string),
			api.DatadogAlertChannelType,
			api.DatadogDataV2{
				DatadogSite: site,
				DatadogType: service,
				ApiKey:      d.Get("api_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		datadog.Enabled = 0
	}

	datadog.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.DatadogAlertChannelType, datadog)
	response, err := lacework.V2.AlertChannels.UpdateDatadog(datadog)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.DatadogAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.DatadogAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.DatadogAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelDatadogDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.DatadogAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.DatadogAlertChannelType, d.Id())
	return nil
}
