package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationDockerHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationDockerHubCreate,
		Read:   resourceLaceworkIntegrationDockerHubRead,
		Update: resourceLaceworkIntegrationDockerHubUpdate,
		Delete: resourceLaceworkIntegrationDockerHubDelete,

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
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"limit_by_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
			},
			"limit_by_label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
			},
			"limit_by_repos": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"limit_num_imgs": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
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
	data := api.NewDockerHubRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:   d.Get("limit_by_tag").(string),
			LimitByLabel: d.Get("limit_by_label").(string),
			LimitByRep:   d.Get("limit_by_repos").(string),
			LimitNumImg:  d.Get("limit_num_imgs").(int),
			Credentials: api.ContainerRegCreds{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), data)
	response, err := lacework.Integrations.CreateContainerRegistry(data)
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.TypeName)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created %s integration %s registry type with guid: %v\n",
			api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), integration.IntgGuid)
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationDockerHubRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), d.Id())
	response, err := lacework.Integrations.GetContainerRegistry(d.Id())

	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)

			d.Set("username", integration.Data.Credentials.Username)
			d.Set("password", integration.Data.Credentials.Password)
			d.Set("limit_by_tag", integration.Data.LimitByTag)
			d.Set("limit_by_label", integration.Data.LimitByLabel)
			d.Set("limit_by_repos", integration.Data.LimitByRep)
			d.Set("limit_num_imgs", integration.Data.LimitNumImg)

			log.Printf("[INFO] Read %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationDockerHubUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	data := api.NewDockerHubRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:   d.Get("limit_by_tag").(string),
			LimitByLabel: d.Get("limit_by_label").(string),
			LimitByRep:   d.Get("limit_by_repos").(string),
			LimitNumImg:  d.Get("limit_num_imgs").(int),
			Credentials: api.ContainerRegCreds{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), data)
	response, err := lacework.Integrations.UpdateContainerRegistry(data)
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)

			log.Printf("[INFO] Updated %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), d.Id())

			return nil
		}
	}

	return nil
}

func resourceLaceworkIntegrationDockerHubDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), d.Id())

	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerHubRegistry.String(), d.Id())

	return nil
}
