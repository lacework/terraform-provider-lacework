package lacework

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationDockerHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationDockerHubCreate,
		Read:   resourceLaceworkIntegrationDockerHubRead,
		Update: resourceLaceworkIntegrationDockerHubUpdate,
		Delete: resourceLaceworkIntegrationDockerHubDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkContainerRegistry,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Docker Hub integration name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Docker user that has at least read-only permissions to the Docker Hub container repositories",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password for the specified Docker Hub user",
			},
			"limit_by_tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A list of image tags to limit the assessment of images with matching tags",
			},
			"limit_by_labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A key based map of labels to limit the assessment of images with matching key:value labels",
			},
			"limit_by_repositories": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A list of repositories to assess",
			},
			"limit_num_imgs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The maximum number of newest container images to assess per repository",
			},
			"non_os_package_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable program language scanning",
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_or_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_level": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkIntegrationDockerHubCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.DockerhubContainerRegistry,
		api.DockerhubData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitByLabel:     castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.DockerhubCredentials{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s registry type with data:\n%+v\n", api.DockerhubContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.Create(data)
	if err != nil {
		return err
	}

	integration := response.Data
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Created %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationDockerHubRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetDockerhub(d.Id())

	if err != nil {
		return resourceNotFound(d, err)
	}

	integration := response.Data
	if integration.IntgGuid == d.Id() {
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)

		d.Set("username", integration.Data.Credentials.Username)
		d.Set("non_os_package_support", integration.Data.NonOSPackageEval)
		d.Set("limit_num_imgs", response.Data.Data.LimitNumImg)
		d.Set("limit_by_tags", response.Data.Data.LimitByTag)
		d.Set("limit_by_repositories", response.Data.Data.LimitByRep)
		d.Set("limit_by_labels", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))

		log.Printf("[INFO] Read %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationDockerHubUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.DockerhubContainerRegistry,
		api.DockerhubData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			LimitByLabel:     castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels"),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.DockerhubCredentials{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s registry type with data:\n%+v\n", api.DockerhubContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateDockerhub(data)
	if err != nil {
		return err
	}

	integration := response.Data
	if integration.IntgGuid == d.Id() {
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Updated %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), d.Id())
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationDockerHubDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), d.Id())

	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s registry type with guid: %v\n", api.DockerhubContainerRegistry.String(), d.Id())

	return nil
}
