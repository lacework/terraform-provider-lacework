package lacework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

func resourceLaceworkPolicyException() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkPolicyExceptionCreate,
		Read:   resourceLaceworkPolicyExceptionRead,
		Update: resourceLaceworkPolicyExceptionUpdate,
		Delete: resourceLaceworkPolicyExceptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkPolicyException,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the policy the exception is associated",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the policy exception",
			},
			"constraint": {
				Type:        schema.TypeSet,
				MinItems:    1,
				Required:    true,
				Description: "The list of constraints",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_key": {
							Type:        schema.TypeString,
							Description: "The field key",
							Required:    true,
						},
						"field_values": {
							Type:        schema.TypeList,
							Description: "The field values",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(val interface{}) string {
									return strings.TrimSpace(val.(string))
								},
							},
						},
						"field_value_map": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A list of key value pairs to filter the policy exception",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "The values map key",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "The values map value",
										Required:    true,
									},
								},
							},
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
		},
	}
}

func resourceLaceworkPolicyExceptionCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		policyID = d.Get("policy_id").(string)
	)

	constraints, err := castSchemaSetToConstraintArray(d, "constraint")
	if err != nil {
		return err
	}
	exception := api.PolicyException{
		Description: d.Get("description").(string),
		Constraints: constraints,
	}

	log.Printf("[INFO] Creating Policy Exception with data:\n%+v\n", exception)
	response, err := lacework.V2.Policy.Exceptions.Create(policyID, exception)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ExceptionID)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)

	log.Printf("[INFO] Created Policy Exception with guid %s\n", response.Data.ExceptionID)
	return nil
}

func resourceLaceworkPolicyExceptionRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		response api.PolicyExceptionResponse
	)

	log.Printf("[INFO] Reading Policy with guid %s\n", d.Id())
	err := lacework.V2.Policy.Exceptions.Get(d.Get("policy_id").(string), d.Id(), &response)
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.ExceptionID)
	d.Set("description", response.Data.Description)
	d.Set("constraint", response.Data.Constraints)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)

	log.Printf("[INFO] Read Policy Exception with guid %s\n", response.Data.ExceptionID)
	return nil
}

func resourceLaceworkPolicyExceptionUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		policyID = d.Get("policy_id").(string)
	)

	constraints, err := castSchemaSetToConstraintArray(d, "constraint")
	if err != nil {
		return err
	}

	exception := api.PolicyException{
		Description: d.Get("description").(string),
		Constraints: constraints,
		ExceptionID: d.Id(),
	}

	log.Printf("[INFO] Updating Policy Exception with data:\n%+v\n", exception)
	response, err := lacework.V2.Policy.Exceptions.Update(policyID, exception)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ExceptionID)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)

	log.Printf("[INFO] Updated Policy Exception with guid %s\n", response.Data.ExceptionID)
	return nil
}

func resourceLaceworkPolicyExceptionDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting Policy with guid %s\n", d.Id())
	err := lacework.V2.Policy.Exceptions.Delete(d.Get("policy_id").(string), d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Policy Exception with guid %s\n", d.Id())
	return nil
}

func importLaceworkPolicyException(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.PolicyExceptionResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Policy Exception with guid: %s\n", d.Id())

	err := lacework.V2.Policy.Exceptions.Get(d.Get("policy_id").(string), d.Id(), &response)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Policy Exception with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Policy Exception found with guid: %s\n", response.Data.ExceptionID)
	return []*schema.ResourceData{d}, nil
}

func castSchemaSetToConstraintArray(d *schema.ResourceData, attr string) (constraints []api.PolicyExceptionConstraint, err error) {
	var (
		list []any
	)

	if d.Get(attr) == nil {
		return nil, errors.Errorf("attribute %s not found", attr)
	}

	list = d.Get(attr).(*schema.Set).List()
	for _, item := range list {
		iMap, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("[WARN] unable to cast constraint %v", item)
			continue
		}

		val := sanitizeConstraintKeys(iMap)
		if _, ok := val["fieldValues"]; ok {
			delete(val, "fieldvalues")
		}

		v, err := json.Marshal(val)
		if err != nil {
			return nil, errors.New("failed to marshall constraint attribute")
		}
		var t api.PolicyExceptionConstraint
		err = json.Unmarshal(v, &t)
		if err != nil {
			return nil, errors.New("failed to unmarshall constraint attribute")
		}
		constraints = append(constraints, t)
	}
	return
}

func sanitizeConstraintKeys(itemMap map[string]any) map[string]any {
	var newMap = make(map[string]any)
	var constraintMapList []any
	for k, v := range itemMap {
		if k == "field_value_map" {
			list := v.(*schema.Set).List()
			if len(list) > 0 {
				constraintMapList = append(constraintMapList, list...)
				newMap["fieldValues"] = constraintMapList
				continue
			}
		}
		newKey := strings.Replace(k, "_", "", -1)
		if newKey != "fieldvaluemap" {
			newMap[newKey] = v
		}
	}
	return newMap
}
