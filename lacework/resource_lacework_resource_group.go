package lacework

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"log"
)

func resourceLaceworkResourceGroup() *schema.Resource {
	filterKey := &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "For fields that support a tag, the key on which to filter.",
	}

	filterValue := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "The values that the predicate should match.",
	}

	filterOperation := &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The operation that should be applied across filters/groups",
	}

	filterField := &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The field on which to apply the predicate.",
	}
	filterFieldName := &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "A custom name for the filter.",
	}

	filterSchema := &schema.Schema{
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

	groupOperator := &schema.Schema{
		Type:        schema.TypeString,
		Description: "The operation to apply (AND/OR)",
		Required:    true,
	}

	groupSchema := &schema.Resource{
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

	return &schema.Resource{
		Create: resourceLaceworkResourceGroupCreate,
		Read:   resourceLaceworkResourceGroupRead,
		Update: resourceLaceworkResourceGroupUpdate,
		Delete: resourceLaceworkResourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	queryJson, err := json.Marshal(rgQuery)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGroupGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("query", response.Data.Query)
	d.Set("group", queryJson)
	d.Set("description", response.Data.Description)
	d.Set("last_updated", response.Data.UpdatedTime)
	d.Set("updated_by", response.Data.UpdatedBy)
	d.Set("type", response.Data.Type)

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

	if response.Data.Query == nil {
		return fmt.Errorf("[ERROR] Resource Group with guid %s not found. "+
			"It either does not exist or is not a V2 Resource Group", d.Id())
	}

	queryJson, err := json.Marshal(response.Data.Query)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGroupGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("query", response.Data.Query)
	d.Set("group", queryJson)
	d.Set("description", response.Data.Description)
	d.Set("last_updated", response.Data.UpdatedTime)
	d.Set("updated_by", response.Data.UpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Read %s Resource Group with guid %s\n",
		response.Data.Type, response.Data.ResourceGroupGuid)
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

	queryJson, err := json.Marshal(data.Query)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGroupGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("query", response.Data.Query)
	d.Set("group", queryJson)
	d.Set("description", response.Data.Description)
	d.Set("last_updated", response.Data.UpdatedTime)
	d.Set("updated_by", response.Data.UpdatedBy)
	d.Set("type", response.Data.Type)

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
