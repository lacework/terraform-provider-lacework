package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelJiraServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelJiraServerCreate,
		Read:   resourceLaceworkAlertChannelJiraServerRead,
		Update: resourceLaceworkAlertChannelJiraServerUpdate,
		Delete: resourceLaceworkAlertChannelJiraServerDelete,

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
			"jira_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issue_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_key": {
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
			"custom_template_file": {
				Type:     schema.TypeString,
				Optional: true,
				// @afiune when we migrate to terraform-plugin-sdk/v2
				//ValidateFunc: validation.StringIsJSON,
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

func resourceLaceworkAlertChannelJiraServerCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraAlertChannelData{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewJiraServerAlertChannel(d.Get("name").(string), jiraData)
	if !d.Get("enabled").(bool) {
		jira.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.JiraIntegration, jira)
	response, err := lacework.Integrations.CreateJiraAlertChannel(jira)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateJiraAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.JiraIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelJiraServerRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	response, err := lacework.Integrations.GetJiraAlertChannel(d.Id())
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

			d.Set("jira_url", integration.Data.JiraUrl)
			d.Set("issue_type", integration.Data.IssueType)
			d.Set("group_issues_by", integration.Data.IssueGrouping)
			d.Set("project_key", integration.Data.ProjectID)
			d.Set("username", integration.Data.Username)

			customTemplateString, err := integration.Data.DecodeCustomTemplateFile()
			if err != nil {
				log.Printf("[ERROR] Unable to decode CustomTemplateFile: %v\n", integration.Data.CustomTemplateFile)
				d.Set("custom_template_file", integration.Data.CustomTemplateFile)
			} else {
				d.Set("custom_template_file", customTemplateString)
			}

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.JiraIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelJiraServerUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraAlertChannelData{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			Password:      d.Get("password").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewJiraServerAlertChannel(d.Get("name").(string), jiraData)
	if !d.Get("enabled").(bool) {
		jira.Enabled = 0
	}

	jira.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.JiraIntegration, jira)
	response, err := lacework.Integrations.UpdateJiraAlertChannel(jira)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateJiraAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelJiraServerDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	return nil
}
