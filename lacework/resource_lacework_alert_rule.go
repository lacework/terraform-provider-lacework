package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lacework/go-sdk/v2/api"
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
			StateContext: importLaceworkAlertRule,
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
			"alert_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: fmt.Sprintf("List of alert categories for the alert rule. Valid categories are: %s",
					strings.Join(api.AlertRuleCategories, ", ")),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
					ValidateFunc: validation.StringInSlice(api.AlertRuleCategories, false),
				},
			},
			"alert_sources": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: fmt.Sprintf("List of alert sources for the alert rule. Valid sources are: %s",
					strings.Join(api.AlertRuleSources, ", ")),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
					ValidateFunc: validation.StringInSlice(api.AlertRuleSources, false),
				},
			},
			"event_categories": {
				Type:       schema.TypeSet,
				Optional:   true,
				Deprecated: "This attribute is deprecated and has been replaced by `alert_subcategories`",
				Description: fmt.Sprintf("List of event categories for the alert rule. Valid categories are: %s",
					strings.Join(api.AlertRuleSubCategories, ", ")),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
					ValidateFunc: validation.StringInSlice(api.AlertRuleSubCategories, false),
				},
			},
			"alert_subcategories": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"event_categories"},
				Description: fmt.Sprintf("List of alert subcategories for the alert rule. Valid categories are: %s",
					strings.Join(api.AlertRuleSubCategories, ", ")),
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
					ValidateFunc: validation.StringInSlice(api.AlertRuleSubCategories, false),
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

	var alertSubcategories []interface{}
	if _, ok := d.GetOk("alert_subcategories"); ok {
		alertSubcategories = d.Get("alert_subcategories").(*schema.Set).List()
	} else if _, ok := d.GetOk("event_categories"); ok {
		alertSubcategories = d.Get("event_categories").(*schema.Set).List()
	}

	var (
		lacework        = meta.(*api.Client)
		resourceGroups  = d.Get("resource_groups").(*schema.Set).List()
		alertCategories = d.Get("alert_categories").(*schema.Set).List()
		alertSources    = d.Get("alert_sources").(*schema.Set).List()
		severities      = api.NewAlertRuleSeverities(castAttributeToStringSlice(d, "severities"))
		alertRule       = api.NewAlertRule(d.Get("name").(string),
			api.AlertRuleConfig{
				Description:        d.Get("description").(string),
				Channels:           castStringSlice(alertChannels),
				Severities:         severities,
				AlertSubCategories: castStringSlice(alertSubcategories),
				AlertCategories:    castStringSlice(alertCategories),
				AlertSources:       castStringSlice(alertSources),
				ResourceGroups:     castStringSlice(resourceGroups),
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
	if _, ok := d.GetOk("alert_subcategories"); ok {
		d.Set("alert_subcategories", response.Data.Filter.AlertSubCategories)
	} else if _, ok := d.GetOk("event_categories"); ok {
		d.Set("event_categories", convertSubCategories(response.Data.Filter.AlertSubCategories))
	}
	d.Set("alert_categories", response.Data.Filter.AlertCategories)
	d.Set("alert_sources", response.Data.Filter.AlertSources)
	d.Set("alert_channels", response.Data.Channels)

	log.Printf("[INFO] Read alert rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var alertChannels []interface{}
	if _, ok := d.GetOk("alert_channels"); ok {
		alertChannels = d.Get("alert_channels").(*schema.Set).List()
	}

	var alertSubcategories []interface{}
	if _, ok := d.GetOk("alert_subcategories"); ok {
		alertSubcategories = d.Get("alert_subcategories").(*schema.Set).List()
	} else if _, ok := d.GetOk("event_categories"); ok {
		alertSubcategories = d.Get("event_categories").(*schema.Set).List()
	}

	var (
		lacework        = meta.(*api.Client)
		resourceGroups  = d.Get("resource_groups").(*schema.Set).List()
		alertCategories = d.Get("alert_categories").(*schema.Set).List()
		alertSources    = d.Get("alert_sources").(*schema.Set).List()
		severities      = api.NewAlertRuleSeverities(castAttributeToStringSlice(d, "severities"))
		alertRule       = api.NewAlertRule(d.Get("name").(string),
			api.AlertRuleConfig{
				Description:        d.Get("description").(string),
				Channels:           castStringSlice(alertChannels),
				Severities:         severities,
				AlertSubCategories: castStringSlice(alertSubcategories),
				AlertCategories:    castStringSlice(alertCategories),
				AlertSources:       castStringSlice(alertSources),
				ResourceGroups:     castStringSlice(resourceGroups),
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

func importLaceworkAlertRule(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

// Convert subCategory values to deprecated eventCatory values
func convertSubCategories(categories []string) []string {
	var res []string
	for _, c := range categories {
		switch c {
		case "Application":
			res = append(res, "App")
		case "Cloud Activity":
			res = append(res, "Cloud")
		case "Kubernetes Activity":
			res = append(res, "K8sActivity")
		default:
			res = append(res, c)
		}
	}
	return res
}
