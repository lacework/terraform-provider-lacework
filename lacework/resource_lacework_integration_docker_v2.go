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
			"non_os_package_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	data := api.NewDockerV2RegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:       limitByTags,
			LimitByLabel:     limitByLabels,
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.ContainerRegCreds{
				Username: d.Get("username").(string),
				Password: d.Get("password").(string),
				SSL:      d.Get("ssl").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), data)
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
			api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), integration.IntgGuid)
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationDockerV2Read(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), d.Id())
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

			d.Set("registry_domain", integration.Data.RegistryDomain)
			d.Set("username", integration.Data.Credentials.Username)
			d.Set("password", integration.Data.Credentials.Password)
			d.Set("ssl", integration.Data.Credentials.SSL)
			d.Set("non_os_package_support", integration.Data.NonOSPackageEval)

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

			log.Printf("[INFO] Read %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationDockerV2Update(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	data := api.NewDockerV2RegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:       limitByTags,
			LimitByLabel:     limitByLabels,
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.ContainerRegCreds{
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

	log.Printf("[INFO] Updating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), data)
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
				api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), d.Id())

			return nil
		}
	}

	return nil
}

func resourceLaceworkIntegrationDockerV2Delete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), d.Id())

	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.DockerV2Registry.String(), d.Id())

	return nil
}
