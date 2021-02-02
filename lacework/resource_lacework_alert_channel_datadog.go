package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"datadog_site": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case api.DatadogSiteEu.String(), api.DatadogSiteCom.String():
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'eu' or 'com'", key,
							),
						}
					}
				},
			},
			"datadog_service": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case api.DatadogServiceLogsDetails.String(), api.DatadogServiceLogsSummary.String(), api.DatadogServiceEventsSummary.String():
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'Logs Detail', 'Logs Summary' or 'Events Summary'", key,
							),
						}
					}
				},
			},
			"api_key": {
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

func resourceLaceworkAlertChannelDatadogCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		datadog  = api.NewDatadogAlertChannel(d.Get("name").(string),
			api.DatadogChannelData{
				DatadogSite:    d.Get("datadog_site").(string),
				DatadogService: d.Get("datadog_service").(string),
				ApiKey:         d.Get("api_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		datadog.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.DatadogChannelIntegration, datadog)
	response, err := lacework.Integrations.CreateDatadogAlertChannel(datadog)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateDatadogAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.DatadogChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelDatadogRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.DatadogChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetDatadogAlertChannel(d.Id())
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
			d.Set("datadog_site", integration.Data.DatadogSite)
			d.Set("datadog_service", integration.Data.DatadogService)
			d.Set("api_key", integration.Data.ApiKey)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.DatadogChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelDatadogUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		datadog  = api.NewDatadogAlertChannel(d.Get("name").(string),
			api.DatadogChannelData{
				DatadogSite:    d.Get("datadog_site").(string),
				DatadogService: d.Get("datadog_service").(string),
				ApiKey:         d.Get("api_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		datadog.Enabled = 0
	}

	datadog.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.DatadogChannelIntegration, datadog)
	response, err := lacework.Integrations.UpdateDatadogAlertChannel(datadog)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateDatadogAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.DatadogChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelDatadogDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.DatadogChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.DatadogChannelIntegration, d.Id())
	return nil
}

func validateDatadogAlertChannelResponse(response *api.DatadogAlertChannelResponse) error {
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
