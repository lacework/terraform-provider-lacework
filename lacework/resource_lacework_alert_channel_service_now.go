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
			"instance_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"issue_grouping": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Events",
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation",
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
		snowData           = api.ServiceNowChannelData{
			InstanceURL:   d.Get("instance_url").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
			IssueGrouping: d.Get("issue_grouping").(string),
		}
		testIntegration = d.Get("test_integration").(bool)
	)

	if len(customTemplateJSON) != 0 {
		snowData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	serviceNow := api.NewServiceNowAlertChannel(d.Get("name").(string), snowData)
	if !d.Get("enabled").(bool) {
		serviceNow.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.ServiceNowChannelIntegration, serviceNow)
	response, err := lacework.Integrations.CreateServiceNowAlertChannel(serviceNow)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateServiceNowAlertChannelResponse(&response)
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

	if testIntegration {
		log.Printf("[INFO] Testing %s integration for guid:%s\n", api.DatadogChannelIntegration, d.Id())
		err := VerifyAlertChannel(d.Id(), lacework)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid: %s successfully \n", api.DatadogChannelIntegration, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.ServiceNowChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelServiceNowRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.ServiceNowChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetServiceNowAlertChannel(d.Id())
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
			d.Set("instance_url", integration.Data.InstanceURL)
			d.Set("username", integration.Data.Username)
			d.Set("issue_grouping", integration.Data.IssueGrouping)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.ServiceNowChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelServiceNowUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		snowData           = api.ServiceNowChannelData{
			InstanceURL:   d.Get("instance_url").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
			IssueGrouping: d.Get("issue_grouping").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		snowData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	serviceNow := api.NewServiceNowAlertChannel(d.Get("name").(string), snowData)
	if !d.Get("enabled").(bool) {
		serviceNow.Enabled = 0
	}

	serviceNow.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.ServiceNowChannelIntegration, serviceNow)
	response, err := lacework.Integrations.UpdateServiceNowAlertChannel(serviceNow)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateServiceNowAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.ServiceNowChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelServiceNowDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.ServiceNowChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.ServiceNowChannelIntegration, d.Id())
	return nil
}

func validateServiceNowAlertChannelResponse(response *api.ServiceNowAlertChannelResponse) error {
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
