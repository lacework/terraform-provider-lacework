package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

var filterKey = &schema.Schema{
	Type:        schema.TypeString,
	Optional:    true,
	Description: "For fields that support a tag, the key on which to filter.",
}

var filterValue = &schema.Schema{
	Type: schema.TypeList,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "The values that the predicate should match.",
}

var filterOperation = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "The operation that should be applied across filters/groups",
}

var filterField = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "The field on which to apply the predicate.",
}

var filterFieldName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "A custom name for the filter.",
}

var filterSchema = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_name": filterFieldName,
			"field":       filterField,
			"operation":   filterOperation,
			"value":       filterValue,
			"key":         filterKey,
		},
	},
}

var groupOperator = &schema.Schema{
	Type:        schema.TypeString,
	Description: "The operation to apply (AND/OR)",
	Required:    true,
}

// Define global variable for groupSchema
var groupSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"operator": groupOperator,
		"filter":   filterSchema,
		"group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"operator": groupOperator,
					"filter":   filterSchema,
					"group": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"operator": groupOperator,
								"filter":   filterSchema,
								"group": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"operator": groupOperator,
											"filter":   filterSchema,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

func resourceLaceworkResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkResourceGroupCreate,
		Read:   resourceLaceworkResourceGroupRead,
		Update: resourceLaceworkResourceGroupUpdate,
		Delete: resourceLaceworkResourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkResourceGroup,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource group name",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the resource group",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the resource group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the resource group",
			},
			"group": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The query used to fetch resources matching the filters and expression",
				Elem:        groupSchema,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group unique identifier",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time in millis when the resource was last updated",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the lacework user who performed the last update",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the resource group is a default resource group.",
			},
		},
	}
}

func importLaceworkResourceGroup(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing resource group.")

	var response api.ResourceGroupResponse
	err := lacework.V2.ResourceGroups.Get(d.Id(), &response)
	if err != nil {
		return nil, err
	}

	err = readResourceGroup(d, &response)

	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func convertRgQueryToInterface(query *api.RGQuery) []interface{} {
	set := []interface{}{}

	result := map[string]interface{}{
		"operator": query.Expression.Operator,
		"filter":   []interface{}{},
		"group":    []interface{}{},
	}

	for _, child := range query.Expression.Children {
		if child.FilterName != "" {
			filter := query.Filters[child.FilterName]
			filterMap := map[string]interface{}{
				"field":       filter.Field,
				"filter_name": child.FilterName,
				"key":         filter.Key, // Assuming key is empty as it's not provided in the context
				"operation":   filter.Operation,
				"value":       filter.Values,
			}
			result["filter"] = append(result["filter"].([]interface{}), filterMap)
		} else {
			nestedGroup := convertRgQueryToInterface(&api.RGQuery{
				Filters:    query.Filters,
				Expression: &api.RGExpression{Children: child.Children, Operator: child.Operator},
			})
			result["group"] = append(result["group"].([]interface{}), nestedGroup...)
		}

		set = append(set, result)
	}

	return set
}

func addFilters(filters *schema.Set, query *api.RGQuery) []string {
	filterList := filters.List()
	var filterNames []string
	for _, f := range filterList {
		filter := f.(map[string]interface{})
		values := filter["value"].([]interface{})
		valuesList := make([]string, len(values))
		for i, v := range values {
			valuesList[i] = fmt.Sprint(v)
		}

		filterName := filter["filter_name"].(string)
		filterNames = append(filterNames, filterName)
		query.Filters[filterName] = &api.RGFilter{
			Field:     filter["field"].(string),
			Operation: filter["operation"].(string),
			Values:    valuesList,
			Key:       filter["key"].(string),
		}
	}

	return filterNames
}

func populateRgQuery(group *schema.Set, query *api.RGQuery, isTopLevelGroup bool) *api.RGChild {
	groupDetails := group.List()
	for _, v := range groupDetails {
		val := v.(map[string]interface{})

		operator := val["operator"].(string)
		nestedGroup := val["group"].(*schema.Set)
		filters := val["filter"].(*schema.Set)
		filterNames := addFilters(filters, query)

		var children []*api.RGChild

		for _, filterName := range filterNames {
			children = append(children, &api.RGChild{
				FilterName: filterName,
			})
		}

		if nestedGroup != nil {
			nestedChild := populateRgQuery(nestedGroup, query, false)
			if nestedChild != nil {
				children = append(children, nestedChild)
			}
		}

		if !isTopLevelGroup {
			return &api.RGChild{
				Children: children,
				Operator: operator,
			}
		}

		rgExpression := &api.RGExpression{
			Operator: operator,
			Children: children,
		}

		query.Expression = rgExpression
	}

	return nil
}

func resourceLaceworkResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	resourceType := d.Get("type").(string)
	groupType, isValid := api.FindResourceGroupType(resourceType)
	if !isValid {
		// This should never reach this. The type is controlled by us in cmd/resource_groups
		return errors.New("internal error")
	}

	rgQuery := api.RGQuery{
		Filters: map[string]*api.RGFilter{},
	}

	populateRgQuery(d.Get("group").(*schema.Set), &rgQuery, true)

	data := api.NewResourceGroup(d.Get("name").(string),
		groupType,
		d.Get("description").(string),
		&rgQuery)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s Resource Group with data:\n%+v\n",
		data.Type, data)
	response, err := lacework.V2.ResourceGroups.Create(data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGroupGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("group", convertRgQueryToInterface(response.Data.Query))
	d.Set("description", response.Data.Description)
	d.Set("last_updated", response.Data.UpdatedTime.UTC().String())
	d.Set("updated_by", response.Data.UpdatedBy)
	d.Set("type", response.Data.Type)
	d.Set("is_default", response.Data.IsDefaultBoolean)

	log.Printf("[INFO] Created %s Resource Group with guid %s\n",
		response.Data.Type, response.Data.ResourceGroupGuid)
	return nil
}

func resourceLaceworkResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading V2 Resource Group with guid %s\n", d.Id())
	var response api.ResourceGroupResponse
	err := lacework.V2.ResourceGroups.Get(d.Id(), &response)
	if err != nil {
		return resourceNotFound(d, err)
	}

	return readResourceGroup(d, &response)
}

func readResourceGroup(d *schema.ResourceData, resourceGroup *api.ResourceGroupResponse) error {
	if resourceGroup.Data.Query == nil {
		return fmt.Errorf("[ERROR] Resource Group with guid %s not found. "+
			"It either does not exist or is not a V2 Resource Group", d.Id())
	}

	d.SetId(resourceGroup.Data.ResourceGroupGuid)
	d.Set("name", resourceGroup.Data.Name)
	d.Set("enabled", resourceGroup.Data.Enabled == 1)
	d.Set("group", convertRgQueryToInterface(resourceGroup.Data.Query))
	d.Set("description", resourceGroup.Data.Description)
	d.Set("last_updated", resourceGroup.Data.UpdatedTime.UTC().String())
	d.Set("updated_by", resourceGroup.Data.UpdatedBy)
	d.Set("type", resourceGroup.Data.Type)
	d.Set("is_default", resourceGroup.Data.IsDefaultBoolean)

	log.Printf("[INFO] Read %s Resource Group with guid %s\n",
		resourceGroup.Data.Type, resourceGroup.Data.ResourceGroupGuid)
	return nil
}

func resourceLaceworkResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	resourceType := d.Get("type").(string)
	groupType, isValid := api.FindResourceGroupType(resourceType)
	if !isValid {
		// This should never reach this. The type is controlled by us in cmd/resource_groups
		return errors.New("internal error")
	}

	rgQuery := api.RGQuery{
		Filters: map[string]*api.RGFilter{},
	}

	populateRgQuery(d.Get("group").(*schema.Set), &rgQuery, true)

	data := api.NewResourceGroup(d.Get("name").(string),
		groupType,
		d.Get("description").(string),
		&rgQuery)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.ResourceGroupGuid = d.Id()

	log.Printf("[INFO] Updating %s Resource Group with data:\n%+v\n",
		data.Type, data)
	response, err := lacework.V2.ResourceGroups.Update(&data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGroupGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("group", convertRgQueryToInterface(response.Data.Query))
	d.Set("description", response.Data.Description)
	d.Set("last_updated", response.Data.UpdatedTime.UTC().String())
	d.Set("updated_by", response.Data.UpdatedBy)
	d.Set("type", response.Data.Type)
	d.Set("is_default", response.Data.IsDefaultBoolean)

	log.Printf("[INFO] Updated %s Resource Group with guid %s\n",
		data.Type, response.Data.ResourceGroupGuid)
	return nil
}

func resourceLaceworkResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting Resource Group with guid %s\n", d.Id())
	err := lacework.V2.ResourceGroups.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Resource Group with guid %s\n", d.Id())
	return nil
}
