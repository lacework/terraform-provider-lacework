package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelJiraCloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelJiraCloudCreate,
		Read:   resourceLaceworkAlertChannelJiraCloudRead,
		Update: resourceLaceworkAlertChannelJiraCloudUpdate,
		Delete: resourceLaceworkAlertChannelJiraCloudDelete,

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
			"api_token": {
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

func resourceLaceworkAlertChannelJiraCloudCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraAlertChannelData{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			ApiToken:      d.Get("api_token").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewJiraCloudAlertChannel(d.Get("name").(string), jiraData)
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

func resourceLaceworkAlertChannelJiraCloudRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceLaceworkAlertChannelJiraCloudUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		customTemplateJSON = d.Get("custom_template_file").(string)
		jiraData           = api.JiraAlertChannelData{
			JiraUrl:       d.Get("jira_url").(string),
			IssueType:     d.Get("issue_type").(string),
			IssueGrouping: d.Get("group_issues_by").(string),
			ProjectID:     d.Get("project_key").(string),
			Username:      d.Get("username").(string),
			ApiToken:      d.Get("api_token").(string),
		}
	)

	if len(customTemplateJSON) != 0 {
		jiraData.EncodeCustomTemplateFile(customTemplateJSON)
	}

	jira := api.NewJiraCloudAlertChannel(d.Get("name").(string), jiraData)
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

func resourceLaceworkAlertChannelJiraCloudDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.JiraIntegration, d.Id())
	return nil
}

// validateJiraAlertChannelResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateJiraAlertChannelResponse(response *api.JiraAlertChannelResponse) error {
	if len(response.Data) == 0 {
		// @afiune this edge case should never happen, if we land here it means that
		// something went wrong in the server side of things (Lacework API), so let
		// us inform that to our users
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		// @afiune if we are creating a single integration and the server returns
		// more than one integration inside the 'data' field, it is definitely another
		// edge case that should never happen
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
