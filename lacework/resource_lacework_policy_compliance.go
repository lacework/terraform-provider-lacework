package lacework

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkPolicyCompliance() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkPolicyComplianceCreate,
		Read:   resourceLaceworkPolicyComplianceRead,
		Update: resourceLaceworkPolicyComplianceUpdate,
		Delete: resourceLaceworkPolicyComplianceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkPolicyCompliance,
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
			"query_language": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: "The language of the query/module",
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
			"policy_id_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string appended to the end of the policy id",
			},
			"remediation": {
				Type:        schema.TypeString,
				Required:    true,
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
			"alerting_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether alerting is enabled or disabled",
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

func resourceLaceworkPolicyComplianceCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	policy := api.NewPolicy{
		PolicyType:   api.PolicyTypeCompliance.String(),
		QueryID:      d.Get("query_id").(string),
		Title:        d.Get("title").(string),
		Enabled:      d.Get("enabled").(bool),
		Description:  d.Get("description").(string),
		Remediation:  d.Get("remediation").(string),
		Severity:     d.Get("severity").(string),
		PolicyID:     d.Get("policy_id_suffix").(string),
		Tags:         castStringSlice(d.Get("tags").(*schema.Set).List()),
		AlertEnabled: d.Get("alerting_enabled").(bool),
	}

	log.Printf("[INFO] Creating Policy with data:\n%+v\n", policy)
	response, err := lacework.V2.Policy.Create(policy)
	if err != nil {
		return err
	}

	d.SetId(response.Data.PolicyID)
	if response.Data.QueryLanguage != nil {
		d.Set("query_language", response.Data.QueryLanguage)
	}
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("computed_tags", strings.Join(response.Data.Tags, ","))

	log.Printf("[INFO] Created Policy with guid %s\n", response.Data.PolicyID)
	return nil
}

func resourceLaceworkPolicyComplianceRead(d *schema.ResourceData, meta interface{}) error {
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
	if response.Data.QueryLanguage != nil {
		d.Set("query_language", response.Data.QueryLanguage)
	}
	d.Set("enabled", response.Data.Enabled)
	d.Set("description", response.Data.Description)
	d.Set("severity", response.Data.Severity)
	d.Set("remediation", response.Data.Remediation)
	d.Set("type", response.Data.PolicyType)
	d.Set("alerting_enabled", response.Data.AlertEnabled)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("computed_tags", strings.Join(response.Data.Tags, ","))

	log.Printf("[INFO] Read Policy with guid %s\n", response.Data.PolicyID)
	return nil
}

func resourceLaceworkPolicyComplianceUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	if d.HasChange("policy_id_suffix") {
		return errors.New("unable to change ID of an existing policy")
	}

	policyEnabled := d.Get("enabled").(bool)

	policy := api.UpdatePolicy{
		PolicyType:  api.PolicyTypeCompliance.String(),
		QueryID:     d.Get("query_id").(string),
		Title:       d.Get("title").(string),
		Enabled:     &policyEnabled,
		Description: d.Get("description").(string),
		Remediation: d.Get("remediation").(string),
		Severity:    d.Get("severity").(string),
		PolicyID:    d.Id(),
		Tags:        castStringSlice(d.Get("tags").(*schema.Set).List()),
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

func resourceLaceworkPolicyComplianceDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting Policy with guid %s\n", d.Id())
	_, err := lacework.V2.Policy.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Policy with guid %s\n", d.Id())
	return nil
}

func importLaceworkPolicyCompliance(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Policy with guid: %s\n", d.Id())

	response, err := lacework.V2.Policy.Get(d.Id())
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Policy with guid '%s' was not found",
			d.Id(),
		)
	}
	if response.Data.PolicyType != api.PolicyTypeCompliance.String() {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Policy with guid '%s' is not a compliance policy",
			d.Id(),
		)
	}
	log.Printf("[INFO] Policy found with guid: %s\n", response.Data.PolicyID)
	return []*schema.ResourceData{d}, nil
}
