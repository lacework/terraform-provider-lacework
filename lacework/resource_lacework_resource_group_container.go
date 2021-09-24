package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"log"
	"strings"
)

func resourceLaceworkResourceGroupContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkResourceGroupContainerCreate,
		Read:   resourceLaceworkResourceGroupContainerRead,
		Update: resourceLaceworkResourceGroupContainerUpdate,
		Delete: resourceLaceworkResourceGroupContainerDelete,

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
			"container_labels": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required:    true,
				Description: "The key value pairs of container labels to include in the resource group",
			},
			"container_tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Required:    true,
				Description: "The list of containers tags to include in the resource group",
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
				Description: "Whether the resource group is a default resource group",
			},
		},
	}
}

func resourceLaceworkResourceGroupContainerCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroup(d.Get("name").(string),
		api.ContainerResourceGroup,
		api.ContainerResourceGroupProps{
			Description:     d.Get("description").(string),
			ContainerLabels: castAttributeToArrayOfKeyValueMap(d, "container_labels"),
			ContainerTags:   castAttributeToStringSlice(d, "container_tags"),
		})

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s Resource Group with data:\n%+v\n",
		api.ContainerResourceGroup.String(), data)
	response, err := lacework.V2.ResourceGroups.CreateContainer(&data)
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
		api.ContainerResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupContainerRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s Resource Group with guid %s\n",
		api.ContainerResourceGroup.String(), d.Id())
	response, err := lacework.V2.ResourceGroups.GetContainer(d.Id())
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
	d.Set("container_tags", response.Data.Props.ContainerTags)
	d.Set("container_labels", response.Data.Props.ContainerLabels)

	log.Printf("[INFO] Read %s Resource Group with guid %s\n",
		api.ContainerResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewResourceGroup(d.Get("name").(string),
		api.ContainerResourceGroup,
		api.ContainerResourceGroupProps{
			Description:     d.Get("description").(string),
			ContainerLabels: castAttributeToArrayOfKeyValueMap(d, "container_labels"),
			ContainerTags:   castAttributeToStringSlice(d, "container_tags"),
		})

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.ResourceGuid = d.Id()

	log.Printf("[INFO] Updating %s Resource Group with data:\n%+v\n",
		api.ContainerResourceGroup.String(), data)
	response, err := lacework.V2.ResourceGroups.UpdateContainer(&data)
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
		api.ContainerResourceGroup.String(), response.Data.ResourceGuid)
	return nil
}

func resourceLaceworkResourceGroupContainerDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s Resource Group with guid %s\n",
		api.ContainerResourceGroup.String(), d.Id())
	err := lacework.V2.ResourceGroups.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s Resource Group with guid %s\n",
		api.ContainerResourceGroup.String(), d.Id())
	return nil
}
