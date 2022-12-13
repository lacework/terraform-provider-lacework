package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelJiraCloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelJiraCloudCreate,
		Read:   resourceLaceworkAlertChannelJiraCloudRead,
		Update: resourceLaceworkAlertChannelJiraCloudUpdate,
		Delete: resourceLaceworkAlertChannelJiraCloudDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAlertChannel,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The alert channel integration name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"jira_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of your Jira implementation without https protocol",
			},
			"issue_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Jira issue type (such as a Bug) to create when a new Jira issue is created",
			},
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project key for the Jira project where the new Jira issues should be created",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Jira user name. Lacework recommends a dedicated Jira user.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The Jira API Token",
			},
			"custom_template_file": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "A Custom Template JSON file to populate fields in the new Jira issues",
			},
			"configuration": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Unidirectional", "Bidirectional":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'Unidirectional' or 'Bidirectional' (default: Unidirectional)", key,
							),
						}
					}
				},
				Description: "Whether the integration is Unidirectional or Bidirectional. Defaults to Unidirectional",
			},
			"group_issues_by": {
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
				Description: "Defines how Lacework compliance events get grouped. Must be one of Events or Resources. Defaults to Events",
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

func resourceLaceworkAlertChannelJiraCloudCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraDataV2{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			Configuration: d.Get("configuration").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			ApiToken:      d.Get("api_token").(string),
			JiraType:      api.JiraCloudAlertType,
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewAlertChannel(d.Get("name").(string), api.JiraAlertChannelType, jiraData)
	if !d.Get("enabled").(bool) {
		jira.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.JiraAlertChannelType, jira)
	response, err := lacework.V2.AlertChannels.Create(jira)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.JiraAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.JiraAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.JiraAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelJiraCloudRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.JiraAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetJira(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)

	d.Set("jira_url", response.Data.Data.JiraUrl)
	d.Set("issue_type", response.Data.Data.IssueType)
	d.Set("configuration", response.Data.Data.Configuration)
	d.Set("group_issues_by", response.Data.Data.IssueGrouping)
	d.Set("project_key", response.Data.Data.ProjectID)
	d.Set("username", response.Data.Data.Username)

	customTemplateString, err := response.Data.Data.DecodeCustomTemplateFile()
	if err != nil {
		log.Printf("[ERROR] Unable to decode CustomTemplateFile: %v\n", response.Data.Data.CustomTemplateFile)
		d.Set("custom_template_file", response.Data.Data.CustomTemplateFile)
	} else {
		d.Set("custom_template_file", customTemplateString)
	}

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.JiraAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelJiraCloudUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraDataV2{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			Configuration: d.Get("configuration").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			ApiToken:      d.Get("api_token").(string),
			JiraType:      api.JiraCloudAlertType,
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewAlertChannel(d.Get("name").(string), api.JiraAlertChannelType, jiraData)
	if !d.Get("enabled").(bool) {
		jira.Enabled = 0
	}

	jira.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.JiraAlertChannelType, jira)
	response, err := lacework.V2.AlertChannels.UpdateJira(jira)
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
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.JiraAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.JiraAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.JiraAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelJiraCloudDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.JiraAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.JiraAlertChannelType, d.Id())
	return nil
}
