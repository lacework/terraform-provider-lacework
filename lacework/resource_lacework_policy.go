package lacework

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkPolicyCreate,
		Read:   resourceLaceworkPolicyRead,
		Update: resourceLaceworkPolicyUpdate,
		Delete: resourceLaceworkPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkPolicy,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the policy",
			},
			"query_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the query",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the query",
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The severity for the policy. Valid severities are: " +
					"Critical, High, Medium, Low, Info",
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace(strings.ToLower(val.(string)))
				},
				ValidateDiagFunc: ValidSeverity(),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy type must be 'Violation'",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Violation":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be 'Violation'", key,
							),
						}
					}
				},
			},
			"evaluation": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Hourly",
				Description: "The evaluation frequency must be either 'Hourly' or 'Daily'",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Hourly", "Daily":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be 'Hourly' or 'Daily'", key,
							),
						}
					}
				},
			},
			"policy_id_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string appended to the end of the policy id",
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntAtMost(5000),
				Description:  "Set the number of records returned by the policy. Maximum value is 5000",
			},
			"remediation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remediation message to display",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the policy",
			},
			"tags": {
				Type:        schema.TypeSet,
				Description: "A list of user specified policy tags",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"alerting": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Alerting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile": {
							Type:        schema.TypeString,
							Description: "The alerting profile id",
							Required:    true,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether alerting is enabled or disabled",
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"computed_tags": {
				Type:        schema.TypeString,
				Description: "All policy tags, server generated and user specified tags",
				Computed:    true,
			},
		},
	}
}

func resourceLaceworkPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	policy := api.NewPolicy{
		PolicyType:    d.Get("type").(string),
		QueryID:       d.Get("query_id").(string),
		Title:         d.Get("title").(string),
		Enabled:       d.Get("enabled").(bool),
		Description:   d.Get("description").(string),
		Remediation:   d.Get("remediation").(string),
		Severity:      d.Get("severity").(string),
		Limit:         d.Get("limit").(int),
		EvalFrequency: d.Get("evaluation").(string),
		AlertEnabled:  d.Get("alerting.0.enabled").(bool),
		AlertProfile:  d.Get("alerting.0.profile").(string),
		PolicyID:      d.Get("policy_id_suffix").(string),
		Tags:          castStringSlice(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[INFO] Creating Policy with data:\n%+v\n", policy)
	response, err := lacework.V2.Policy.Create(policy)
	if err != nil {
		return err
	}

	d.SetId(response.Data.PolicyID)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("computed_tags", strings.Join(response.Data.Tags, ","))

	log.Printf("[INFO] Created Policy with guid %s\n", response.Data.PolicyID)
	return nil
}

func resourceLaceworkPolicyRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	log.Printf("[INFO] Reading Policy with guid %s\n", d.Id())
	response, err := lacework.V2.Policy.Get(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.PolicyID)
	d.Set("title", response.Data.Title)
	d.Set("query_id", response.Data.QueryID)
	d.Set("enabled", response.Data.Enabled)
	d.Set("description", response.Data.Description)
	d.Set("evaluation", response.Data.EvalFrequency)
	d.Set("severity", response.Data.Severity)
	d.Set("remediation", response.Data.Remediation)
	d.Set("limit", response.Data.Limit)
	d.Set("type", response.Data.PolicyType)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("computed_tags", strings.Join(response.Data.Tags, ","))

	alerting := make(map[string]interface{})
	alerting["enabled"] = response.Data.AlertEnabled
	alerting["profile"] = response.Data.AlertProfile
	d.Set("alerting", alerting)

	log.Printf("[INFO] Read Policy with guid %s\n", response.Data.PolicyID)
	return nil
}

func resourceLaceworkPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	if d.HasChange("policy_id_suffix") {
		return errors.New("unable to change ID of an existing policy")
	}

	policyEnabled := d.Get("enabled").(bool)
	alertingEnabled := d.Get("alerting.0.enabled").(bool)
	policyLimit := d.Get("limit").(int)

	policy := api.UpdatePolicy{
		PolicyType:    d.Get("type").(string),
		QueryID:       d.Get("query_id").(string),
		Title:         d.Get("title").(string),
		Enabled:       &policyEnabled,
		Description:   d.Get("description").(string),
		Remediation:   d.Get("remediation").(string),
		Severity:      d.Get("severity").(string),
		Limit:         &policyLimit,
		EvalFrequency: d.Get("evaluation").(string),
		AlertEnabled:  &alertingEnabled,
		AlertProfile:  d.Get("alerting.0.profile").(string),
		PolicyID:      d.Id(),
		Tags:          castStringSlice(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[INFO] Updating Policy with data:\n%+v\n", policy)
	response, err := lacework.V2.Policy.Update(policy)
	if err != nil {
		return err
	}

	d.SetId(response.Data.PolicyID)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("computed_tags", strings.Join(response.Data.Tags, ","))

	log.Printf("[INFO] Updated Policy with guid %s\n", response.Data.PolicyID)
	return nil
}

func resourceLaceworkPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting Policy with guid %s\n", d.Id())
	_, err := lacework.V2.Policy.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Policy with guid %s\n", d.Id())
	return nil
}

func importLaceworkPolicy(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Policy with guid: %s\n", d.Id())

	response, err := lacework.V2.Policy.Get(d.Id())
	if err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Policy with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Policy found with guid: %s\n", response.Data.PolicyID)
	return []*schema.ResourceData{d}, nil
}
