package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelQRadar() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelQRadarCreate,
		Read:   resourceLaceworkAlertChannelQRadarRead,
		Update: resourceLaceworkAlertChannelQRadarUpdate,
		Delete: resourceLaceworkAlertChannelQRadarDelete,

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
			"communication_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(api.QRadarCommHttps),
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case string(api.QRadarCommHttps), string(api.QRadarCommHttpsSelfSigned):
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either '%s' or '%s'", key,
								api.QRadarCommHttps, api.QRadarCommHttpsSelfSigned,
							),
						}
					}
				},
			},
			"host_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 65536 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 65535 inclusive, got: %d", key, v))
					}
					return
				},
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

func resourceLaceworkAlertChannelQRadarCreate(d *schema.ResourceData, meta interface{}) error {
	comm, _ := api.QRadarComm(d.Get("communication_type").(string))

	var (
		lacework = meta.(*api.Client)
		qradar   = api.NewQRadarAlertChannel(d.Get("name").(string),
			api.QRadarChannelData{
				HostURL:           d.Get("host_url").(string),
				HostPort:          d.Get("host_port").(int),
				CommunicationType: comm,
			},
		)
	)
	if !d.Get("enabled").(bool) {
		qradar.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.QRadarChannelIntegration, qradar)
	response, err := lacework.Integrations.CreateQRadarAlertChannel(qradar)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateQRadarAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.QRadarChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelQRadarRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.QRadarChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetQRadarAlertChannel(d.Id())
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
			d.Set("host_url", integration.Data.HostURL)
			d.Set("host_port", integration.Data.HostPort)
			d.Set("communicaton_type", integration.Data.CommunicationType)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.QRadarChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelQRadarUpdate(d *schema.ResourceData, meta interface{}) error {
	comm, _ := api.QRadarComm(d.Get("communication_type").(string))

	var (
		lacework = meta.(*api.Client)
		qradar   = api.NewQRadarAlertChannel(d.Get("name").(string),
			api.QRadarChannelData{
				HostURL:           d.Get("host_url").(string),
				HostPort:          d.Get("host_port").(int),
				CommunicationType: comm,
			},
		)
	)

	if !d.Get("enabled").(bool) {
		qradar.Enabled = 0
	}

	qradar.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.QRadarChannelIntegration, qradar)
	response, err := lacework.Integrations.UpdateQRadarAlertChannel(qradar)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateQRadarAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.QRadarChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelQRadarDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.QRadarChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.QRadarChannelIntegration, d.Id())
	return nil
}

func validateQRadarAlertChannelResponse(response *api.QRadarAlertChannelResponse) error {
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
