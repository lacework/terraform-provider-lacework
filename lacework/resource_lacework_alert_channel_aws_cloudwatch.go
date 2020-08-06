package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelAwsCloudWatch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelAwsCloudWatchCreate,
		Read:   resourceLaceworkAlertChannelAwsCloudWatchRead,
		Update: resourceLaceworkAlertChannelAwsCloudWatchUpdate,
		Delete: resourceLaceworkAlertChannelAwsCloudWatchDelete,

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
			"event_bus_arn": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceLaceworkAlertChannelAwsCloudWatchCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAwsCloudWatchAlertChannel(d.Get("name").(string),
			api.AwsCloudWatchData{
				EventBusArn:   d.Get("event_bus_arn").(string),
				IssueGrouping: d.Get("group_issues_by").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AwsCloudWatchIntegration, alert)
	response, err := lacework.Integrations.CreateAwsCloudWatchAlertChannel(alert)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsCloudWatchAlertChannelResponse(&response)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.AwsCloudWatchIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.AwsCloudWatchIntegration, d.Id())
	response, err := lacework.Integrations.GetAwsCloudWatchAlertChannel(d.Id())
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

			d.Set("event_bus_arn", integration.Data.EventBusArn)
			d.Set("group_issues_by", integration.Data.IssueGrouping)

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.AwsCloudWatchIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alert    = api.NewAwsCloudWatchAlertChannel(d.Get("name").(string),
			api.AwsCloudWatchData{
				EventBusArn:   d.Get("event_bus_arn").(string),
				IssueGrouping: d.Get("group_issues_by").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		alert.Enabled = 0
	}

	alert.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsCloudWatchIntegration, alert)
	response, err := lacework.Integrations.UpdateAwsCloudWatchAlertChannel(alert)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsCloudWatchAlertChannelResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.AwsCloudWatchIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelAwsCloudWatchDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.AwsCloudWatchIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.AwsCloudWatchIntegration, d.Id())
	return nil
}

// validateAwsCloudWatchAlertChannelResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateAwsCloudWatchAlertChannelResponse(response *api.AwsCloudWatchResponse) error {
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
