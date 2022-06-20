package lacework

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkPolicyException() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkPolicyExceptionCreate,
		Read:   resourceLaceworkPolicyExceptionRead,
		Update: resourceLaceworkPolicyExceptionUpdate,
		Delete: resourceLaceworkPolicyExceptionDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkPolicy,
		},

		Schema: map[string]*schema.Schema{
			"policyID": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the policy the exception is associated with",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the query",
			},
			"constraints": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the query",
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
		},
	}
}

func resourceLaceworkPolicyExceptionCreate(d *schema.ResourceData, meta interface{}) error {
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

func importLaceworkPolicy(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
