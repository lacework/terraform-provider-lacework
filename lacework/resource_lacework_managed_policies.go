package lacework

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkManagedPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkManagedPoliciesCreate,
		Update: resourceLaceworkManagedPoliciesUpdate,
		Delete: resourceLaceworkManagedPoliciesDelete,
		Read:   resourceLaceworkManagedPoliciesRead,

		Importer: &schema.ResourceImporter{
			State: importLaceworkManagedPolicies,
		},

		Schema: map[string]*schema.Schema{
			"policy": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The id of the policy",
							Required:    true,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "The state of the policy",
						},
						"severity": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "The severity for the policy. Valid severities are: " +
								"Critical, High, Medium, Low, Info",
							StateFunc: func(val interface{}) string {
								return strings.TrimSpace(strings.ToLower(val.(string)))
							},
							ValidateDiagFunc: ValidSeverity(),
						},
					},
				},
				Required:    true,
				Description: "A list of Lacework managed policies",
			},
		},
	}
}

func resourceLaceworkManagedPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(time.Now().UTC().String())
	return resourceLaceworkManagedPoliciesUpdate(d, meta)
}

func resourceLaceworkManagedPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	policies, err := getBulkUpdatePolicies(d, meta)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Policies with data:\n%+v\n", policies)
	_, updateErr := lacework.V2.Policy.UpdateMany(policies)
	if updateErr != nil {
		return updateErr
	}
	log.Printf("[INFO] Updated Policies with data:\n%+v\n", policies)
	return nil
}

func resourceLaceworkManagedPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLaceworkManagedPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func importLaceworkManagedPolicies(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Print("[INFO] Importing Lacework managed policies")
	policies, err := getBulkUpdatePolicies(d, meta)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Found Lacework managed policies: \n%+v\n", policies)
	return []*schema.ResourceData{d}, nil
}

func getBulkUpdatePolicies(d *schema.ResourceData, meta interface{}) (api.BulkUpdatePolicies, error) {
	var policies api.BulkUpdatePolicies
	list := d.Get("policy").(*schema.Set).List()
	seen := make(map[string]bool, 0)

	for _, v := range list {
		val := v.(map[string]interface{})

		if val["id"] == nil || val["id"] == "" {
			continue
		}

		policyID := val["id"].(string)
		enabled := val["enabled"].(bool)

		if !strings.HasPrefix(policyID, "lacework-global") {
			return nil, fmt.Errorf("Unable to update custom policy ID: %s", policyID)
		}
		if seen[policyID] == true {
			return nil, fmt.Errorf("Unable to update duplicate policy ID: %s", policyID)
		}

		policy := api.BulkUpdatePolicy{
			PolicyID: policyID,
			Enabled:  &enabled,
		}

		if val["severity"] != nil && val["severity"] != "" {
			severity := val["severity"].(string)
			policy.Severity = severity
		}

		seen[policyID] = true
		policies = append(policies, policy)
	}
	return policies, nil
}
