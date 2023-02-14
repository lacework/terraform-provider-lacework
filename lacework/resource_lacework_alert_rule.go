package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func resourceLaceworkAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertRuleCreate,
		Read:   resourceLaceworkAlertRuleRead,
		Update: resourceLaceworkAlertRuleUpdate,
		Delete: resourceLaceworkAlertRuleDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAlertRule,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the alert rule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the alert rule",
			},
			"alert_channels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of alert channels for the alert rule",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
			},
			"severities": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Description: "List of severities for the alert rule. Valid severities are:" +
					" Critical, High, Medium, Low, Info",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(cases.Title(language.English).String(strings.ToLower(val.(string))))
					},
					ValidateFunc: func(value interface{}, key string) ([]string, []error) {
						switch strings.ToLower(value.(string)) {
						case "critical", "high", "medium", "low", "info":
							return nil, nil
						default:
							return nil, []error{
								fmt.Errorf(
									"%s: can only be 'Critical', 'High', 'Medium', 'Low', 'Info'", key,
								),
							}
						}
					},
				},
			},
			"resource_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of resource groups for the alert rule",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
			},
			"event_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "List of event categories for the alert rule. Valid categories are: " +
					"Compliance, App, Cloud, File, Machine, User, Platform, K8sActivity",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
					ValidateFunc: func(value interface{}, key string) ([]string, []error) {
						switch value.(string) {
						case "Compliance", "App", "Cloud", "File", "Machine", "User", "Platform", "K8sActivity":
							return nil, nil
						default:
							return nil, []error{
								fmt.Errorf(
									"%s: can only be 'Compliance', 'App', 'Cloud', 'File', 'Machine', 'User', 'Platform', 'K8sActivity'", key,
								),
							}
						}
					},
				},
			},
			"guid": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var alertChannels []interface{}
	if _, ok := d.GetOk("alert_channels"); ok {
		alertChannels = d.Get("alert_channels").(*schema.Set).List()
	}

	var (
		lacework        = meta.(*api.Client)
		resourceGroups  = d.Get("resource_groups").(*schema.Set).List()
		eventCategories = d.Get("event_categories").(*schema.Set).List()
		severities      = api.NewAlertRuleSeverities(castAttributeToStringSlice(d, "severities"))
		alertRule       = api.NewAlertRule(d.Get("name").(string),
			api.AlertRuleConfig{
				Description:     d.Get("description").(string),
				Channels:        castStringSlice(alertChannels),
				Severities:      severities,
				EventCategories: castStringSlice(eventCategories),
				ResourceGroups:  castStringSlice(resourceGroups),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		alertRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Creating alert rule with data:\n%+v\n", alertRule)
	response, err := lacework.V2.AlertRules.Create(alertRule)
	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("guid", response.Data.Guid)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Created alert rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		response api.AlertRuleResponse
	)

	log.Printf("[INFO] Reading alert rule with guid %s\n", d.Id())
	err := lacework.V2.AlertRules.Get(d.Id(), &response)
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("guid", response.Data.Guid)
	d.Set("description", response.Data.Filter.Description)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("severities", api.NewAlertRuleSeveritiesFromIntSlice(response.Data.Filter.Severity).ToStringSlice())
	d.Set("resource_groups", response.Data.Filter.ResourceGroups)
	d.Set("event_categories", response.Data.Filter.EventCategories)
	d.Set("alert_channels", response.Data.Channels)

	log.Printf("[INFO] Read alert rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var alertChannels []interface{}
	if _, ok := d.GetOk("alert_channels"); ok {
		alertChannels = d.Get("alert_channels").(*schema.Set).List()
	}

	var (
		lacework        = meta.(*api.Client)
		resourceGroups  = d.Get("resource_groups").(*schema.Set).List()
		eventCategories = d.Get("event_categories").(*schema.Set).List()
		severities      = api.NewAlertRuleSeverities(castAttributeToStringSlice(d, "severities"))
		alertRule       = api.NewAlertRule(d.Get("name").(string),
			api.AlertRuleConfig{
				Description:     d.Get("description").(string),
				Channels:        castStringSlice(alertChannels),
				Severities:      severities,
				EventCategories: castStringSlice(eventCategories),
				ResourceGroups:  castStringSlice(resourceGroups),
			},
		)
	)

	alertRule.Guid = d.Id()

	if !d.Get("enabled").(bool) {
		alertRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Updating alert rule with data:\n%+v\n", alertRule)
	response, err := lacework.V2.AlertRules.Update(alertRule)

	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("description", response.Data.Filter.Description)
	d.Set("guid", response.Data.Guid)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Updated alert rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting alert rule with guid %s\n", d.Id())
	err := lacework.V2.AlertRules.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted alert rule with guid %s\n", d.Id())
	return nil
}

func importLaceworkAlertRule(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.AlertRuleResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Alert Rule with guid: %s\n", d.Id())

	if err := lacework.V2.AlertRules.Get(d.Id(), &response); err != nil {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Alert Rule with guid '%s' was not found.",
			d.Id(),
		)
	}
	log.Printf("[INFO] Alert Rule found with guid: %s\n", response.Data.Guid)
	return []*schema.ResourceData{d}, nil
}
