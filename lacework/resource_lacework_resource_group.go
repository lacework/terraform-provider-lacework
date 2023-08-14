package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkResourceGroup() *schema.Resource {
	filterKey := &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "For fields that support a tag, the key on which to filter.",
	}

	filterValue := &schema.Schema{
		Elem: &schema.Schema{
			Type: schema.TypeString,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
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
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"filterName": filterFieldName,
				"field":      filterField,
				"operation":  filterOperation,
				"value":      filterValue,
				"key":        filterKey,
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
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": groupOperator,
						"filter":   filterSchema,
						"group": {
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": groupOperator,
									"filter":   filterSchema,
									"group": {
										Type: schema.TypeSet,
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
			"lacework_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lacework account id",
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the resource group",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the resource group is a default resource group.",
			},
		},
	}
}

func resourceLaceworkResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroupWithQuery(d.Get("name").(string),
		d.Get("type").(api.ResourceGroupType),
		d.Get("description").(string),
		d.Get("query").(*api.RGQuery))

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s Resource Group with data:\n%+v\n",
		data.Type, data)
	response, err := lacework.V2.ResourceGroups.Create(&data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGuid)
	d.Set("name", response.Data.NameV2)
	d.Set("lacework_account_id", response.Data.Guid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("query", response.Data.Query)
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

	log.Printf("[INFO] Reading %s Resource Group with guid %s\n",
		api.AwsResourceGroup.String(), d.Id())
	response, err := lacework.V2.ResourceGroups.GetAws(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.ResourceGuid)
	d.Set("name", response.Data.Name)
	d.Set("lacework_account_id", response.Data.Guid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("description", response.Data.Props.Description)
	d.Set("last_updated", response.Data.Props.LastUpdated)
	d.Set("updated_by", response.Data.Props.UpdatedBy)
	d.Set("type", response.Data.Type)
	d.Set("accounts", response.Data.Props.AccountIDs)

	log.Printf("[INFO] Read %s Resource Group with guid %s\n",
		api.AwsResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroup(d.Get("name").(string),
		api.AwsResourceGroup,
		api.AwsResourceGroupProps{
			Description: d.Get("description").(string),
			AccountIDs:  castAttributeToStringSlice(d, "accounts"),
		})

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.ResourceGuid = d.Id()

	log.Printf("[INFO] Updating %s Resource Group with data:\n%+v\n",
		api.AwsResourceGroup.String(), data)
	response, err := lacework.V2.ResourceGroups.UpdateAws(&data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGuid)
	d.Set("name", response.Data.Name)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("last_updated", response.Data.Props.LastUpdated)
	d.Set("updated_by", response.Data.Props.UpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Updated %s Resource Group with guid %s\n",
		api.AwsResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s Resource Group with guid %s\n",
		api.AwsResourceGroup.String(), d.Id())
	err := lacework.V2.ResourceGroups.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s Resource Group with guid %s\n",
		api.AwsResourceGroup.String(), d.Id())
	return nil
}
