package lacework

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationDockerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationDockerV2Create,
		Read:   resourceLaceworkIntegrationDockerV2Read,
		Update: resourceLaceworkIntegrationDockerV2Update,
		Delete: resourceLaceworkIntegrationDockerV2Delete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"registry_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notifications": {
				Type:        schema.TypeBool,
				Description: "Subscribe to registry notifications",
				Optional:    true,
				Default:     false,
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

func resourceLaceworkIntegrationDockerV2Create(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	notifications := d.Get("notifications").(bool)
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.DockerhubV2ContainerRegistry,
		api.DockerhubV2Data{
			LimitByTag:            castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:          castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels"),
			RegistryDomain:        d.Get("registry_domain").(string),
			NonOSPackageEval:      d.Get("non_os_package_support").(bool),
			RegistryNotifications: &notifications,
			Credentials: api.DockerhubV2Credentials{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
				SSL:      d.Get("ssl").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s registry type with data:\n%+v\n", api.DockerhubV2ContainerRegistry.String(), data)
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

	log.Printf("[INFO] Created %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationDockerV2Read(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetDockerhubV2(d.Id())

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

		d.Set("registry_domain", integration.Data.RegistryDomain)
		d.Set("username", integration.Data.Credentials.Username)
		d.Set("password", integration.Data.Credentials.Password)
		d.Set("ssl", integration.Data.Credentials.SSL)
		d.Set("non_os_package_support", integration.Data.NonOSPackageEval)
		d.Set("notifications", integration.Data.RegistryNotifications)
		d.Set("limit_by_tags", response.Data.Data.LimitByTag)
		d.Set("limit_by_labels", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))

		log.Printf("[INFO] Read %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationDockerV2Update(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	notifications := d.Get("notifications").(bool)
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.DockerhubV2ContainerRegistry,
		api.DockerhubV2Data{
			LimitByTag:            castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:          castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels"),
			RegistryDomain:        d.Get("registry_domain").(string),
			NonOSPackageEval:      d.Get("non_os_package_support").(bool),
			RegistryNotifications: &notifications,
			Credentials: api.DockerhubV2Credentials{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
				SSL:      d.Get("ssl").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s registry type with data:\n%+v\n", api.DockerhubV2ContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateDockerhubV2(data)
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

		log.Printf("[INFO] Updated %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), d.Id())

		return nil
	}

	return nil
}

func resourceLaceworkIntegrationDockerV2Delete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), d.Id())

	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s registry type with guid: %v\n", api.DockerhubV2ContainerRegistry.String(), d.Id())

	return nil
}
