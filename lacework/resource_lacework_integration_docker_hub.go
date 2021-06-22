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
			State: importLaceworkIntegration,
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

			// TODO @afiune remove these resources when we release v1.0
			"limit_by_tag": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "*",
				Description:   "A comma-separated list of image tags to limit the assessment of images with matching tags",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_tags` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_tags"},
			},
			"limit_by_label": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "*",
				Description:   "A comma-separated list of image labels to limit the assessment of images with matching labels",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_labels` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_labels"},
			},

			"limit_by_repos": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "A comma-separated list of repositories to assess",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_repositories` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_repositories"},
			},
			// END TODO @afiune

			"limit_by_tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A list of image tags to limit the assessment of images with matching tags",
				ConflictsWith: []string{"limit_by_tag"},
			},
			"limit_by_labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A key based map of labels to limit the assessment of images with matching key:value labels",
				ConflictsWith: []string{"limit_by_label"},
			},
			"limit_by_repositories": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A list of repositories to assess",
				ConflictsWith: []string{"limit_by_repos"},
			},
			"limit_num_imgs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The maximum number of newest container images to assess per repository",
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

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	limitByRepos := d.Get("limit_by_repos").(string)
	if repos := castAttributeToStringSlice(d, "limit_by_repositories"); len(repos) != 0 {
		limitByRepos = strings.Join(repos, ",")
	}

	data := api.NewDockerHubRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:   limitByTags,
			LimitByLabel: limitByLabels,
			LimitByRep:   limitByRepos,
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
			d.Set("limit_num_imgs", integration.Data.LimitNumImg)

			if _, ok := d.GetOk("limit_by_tags"); ok {
				d.Set("limit_by_tags", strings.Split(integration.Data.LimitByTag, ","))
			} else {
				d.Set("limit_by_tag", integration.Data.LimitByTag)
			}

			if _, ok := d.GetOk("limit_by_labels"); ok {
				d.Set("limit_by_labels", strings.Split(integration.Data.LimitByLabel, ","))
			} else {
				d.Set("limit_by_label", integration.Data.LimitByLabel)
			}

			if _, ok := d.GetOk("limit_by_repositories"); ok {
				d.Set("limit_by_repositories", strings.Split(integration.Data.LimitByRep, ","))
			} else {
				d.Set("limit_by_repos", integration.Data.LimitByRep)
			}

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

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	limitByRepos := d.Get("limit_by_repos").(string)
	if repos := castAttributeToStringSlice(d, "limit_by_repositories"); len(repos) != 0 {
		limitByRepos = strings.Join(repos, ",")
	}

	data := api.NewDockerHubRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:   limitByTags,
			LimitByLabel: limitByLabels,
			LimitByRep:   limitByRepos,
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
