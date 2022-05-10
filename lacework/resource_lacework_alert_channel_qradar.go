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
			"communication_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The communication protocol used",
				Default:     string(api.QRadarCommHttps),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain name or IP address of QRadar",
			},
			"host_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The listen port defined in QRadar",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 0 || v > 65536 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 65535 inclusive, got: %d", key, v))
					}
					return
				},
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

func resourceLaceworkAlertChannelQRadarCreate(d *schema.ResourceData, meta interface{}) error {
	comm, _ := api.QRadarComm(d.Get("communication_type").(string))

	var (
		lacework = meta.(*api.Client)
		qradar   = api.NewAlertChannel(d.Get("name").(string),
			api.IbmQRadarAlertChannelType,
			api.IbmQRadarDataV2{
				HostURL:        d.Get("host_url").(string),
				HostPort:       d.Get("host_port").(int),
				QRadarCommType: comm,
			},
		)
	)
	if !d.Get("enabled").(bool) {
		qradar.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.IbmQRadarAlertChannelType, qradar)
	response, err := lacework.V2.AlertChannels.Create(qradar)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
		err := VerifyAlertChannelAndRollback(d, lacework)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.IbmQRadarAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.IbmQRadarAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelQRadarRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetIbmQRadar(d.Id())
	if err != nil {
		return resourceNotFound(d, err, d.Id())
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)
	d.Set("host_url", integration.Data.HostURL)
	d.Set("host_port", integration.Data.HostPort)
	d.Set("communicaton_type", integration.Data.QRadarCommType)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.IbmQRadarAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelQRadarUpdate(d *schema.ResourceData, meta interface{}) error {
	comm, _ := api.QRadarComm(d.Get("communication_type").(string))

	var (
		lacework = meta.(*api.Client)
		qradar   = api.NewAlertChannel(d.Get("name").(string),
			api.IbmQRadarAlertChannelType,
			api.IbmQRadarDataV2{
				HostURL:        d.Get("host_url").(string),
				HostPort:       d.Get("host_port").(int),
				QRadarCommType: comm,
			},
		)
	)

	if !d.Get("enabled").(bool) {
		qradar.Enabled = 0
	}

	qradar.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.IbmQRadarAlertChannelType, qradar)
	response, err := lacework.V2.AlertChannels.UpdateIbmQRadar(qradar)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
		err := lacework.V2.AlertChannels.Test(d.Id())
		if err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.IbmQRadarAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelQRadarDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.IbmQRadarAlertChannelType, d.Id())
	return nil
}
