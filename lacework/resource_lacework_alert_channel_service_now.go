package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelServiceNow() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelServiceNowCreate,
		Read:   resourceLaceworkAlertChannelServiceNowRead,
		Update: resourceLaceworkAlertChannelServiceNowUpdate,
		Delete: resourceLaceworkAlertChannelServiceNowDelete,

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
			"instance_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ServiceNow instance URL",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ServiceNow username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The ServiceNow password",
			},
			"issue_grouping": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Events",
				Description: "Defines how Lacework compliance events get grouped. Must be one of Events or Resources",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Events", "Resources":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'Events' or 'Resources' (default: Events)", key,
							),
						}
					}
				},
			},
			"custom_template_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Populate fields in the ServiceNow incident with values from a custom template JSON file",
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

func resourceLaceworkAlertChannelServiceNowCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		snowData           = api.ServiceNowRestDataV2{
			InstanceURL:   d.Get("instance_url").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
			IssueGrouping: d.Get("issue_grouping").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		snowData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	serviceNow := api.NewAlertChannel(d.Get("name").(string),
		api.ServiceNowRestAlertChannelType,
		snowData)
	if !d.Get("enabled").(bool) {
		serviceNow.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.ServiceNowRestAlertChannelType, serviceNow)
	response, err := lacework.V2.AlertChannels.Create(serviceNow)
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

	log.Printf("[INFO] Created %s integration with guid %s\n", api.ServiceNowRestAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelServiceNowRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.ServiceNowRestAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetServiceNowRest(d.Id())
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
	d.Set("instance_url", integration.Data.InstanceURL)
	d.Set("username", integration.Data.Username)
	d.Set("issue_grouping", integration.Data.IssueGrouping)

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.ServiceNowRestAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelServiceNowUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		snowData           = api.ServiceNowRestDataV2{
			InstanceURL:   d.Get("instance_url").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
			IssueGrouping: d.Get("issue_grouping").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		snowData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	serviceNow := api.NewAlertChannel(d.Get("name").(string),
		api.ServiceNowRestAlertChannelType,
		snowData)
	if !d.Get("enabled").(bool) {
		serviceNow.Enabled = 0
	}

	serviceNow.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.ServiceNowRestAlertChannelType, serviceNow)
	response, err := lacework.V2.AlertChannels.UpdateServiceNowRest(serviceNow)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.ServiceNowRestAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.ServiceNowRestAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.ServiceNowRestAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelServiceNowDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.ServiceNowRestAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.ServiceNowRestAlertChannelType, d.Id())
	return nil
}
