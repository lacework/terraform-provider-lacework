package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkResourceGroupGcp() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkResourceGroupGcpCreate,
		Read:   resourceLaceworkResourceGroupGcpRead,
		Update: resourceLaceworkResourceGroupGcpUpdate,
		Delete: resourceLaceworkResourceGroupGcpDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkResourceGroup,
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
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GCP organization id",
			},
			"projects": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Required:    true,
				Description: "The list of GCP project id's to include in the resource group",
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

func resourceLaceworkResourceGroupGcpCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroup(d.Get("name").(string),
		api.GcpResourceGroup,
		api.GcpResourceGroupProps{
			Description:  d.Get("description").(string),
			Organization: d.Get("organization").(string),
			Projects:     castAttributeToStringSlice(d, "projects"),
		})

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s Resource Group with data:\n%+v\n",
		api.GcpResourceGroup.String(), data)
	response, err := lacework.V2.ResourceGroups.CreateGcp(&data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGuid)
	d.Set("name", response.Data.Name)
	d.Set("lacework_account_id", response.Data.Guid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("description", response.Data.Props.Description)
	d.Set("last_updated", response.Data.Props.LastUpdated)
	d.Set("updated_by", response.Data.Props.UpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Created %s Resource Group with guid %s\n",
		api.GcpResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupGcpRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s Resource Group with guid %s\n",
		api.GcpResourceGroup.String(), d.Id())
	response, err := lacework.V2.ResourceGroups.GetGcp(d.Id())
	if err != nil {
		return err
	}

	d.SetId(response.Data.ResourceGuid)
	d.Set("name", response.Data.Name)
	d.Set("lacework_account_id", response.Data.Guid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("description", response.Data.Props.Description)
	d.Set("last_updated", response.Data.Props.LastUpdated)
	d.Set("updated_by", response.Data.Props.UpdatedBy)
	d.Set("type", response.Data.Type)
	d.Set("projects", response.Data.Props.Projects)
	d.Set("organization", response.Data.Props.Organization)

	log.Printf("[INFO] Read %s Resource Group with guid %s\n",
		api.GcpResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupGcpUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroup(d.Get("name").(string),
		api.GcpResourceGroup,
		api.GcpResourceGroupProps{
			Description:  d.Get("description").(string),
			Organization: d.Get("organization").(string),
			Projects:     castAttributeToStringSlice(d, "projects"),
		})

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.ResourceGuid = d.Id()

	log.Printf("[INFO] Updating %s Resource Group with data:\n%+v\n",
		api.GcpResourceGroup.String(), data)
	response, err := lacework.V2.ResourceGroups.UpdateGcp(&data)
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
		api.GcpResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupGcpDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s Resource Group with guid %s\n",
		api.GcpResourceGroup.String(), d.Id())
	err := lacework.V2.ResourceGroups.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s Resource Group with guid %s\n",
		api.GcpResourceGroup.String(), d.Id())
	return nil
}
