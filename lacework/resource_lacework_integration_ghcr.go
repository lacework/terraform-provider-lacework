package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationGhcr() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGhcrCreate,
		Read:   resourceLaceworkIntegrationGhcrRead,
		Update: resourceLaceworkIntegrationGhcrUpdate,
		Delete: resourceLaceworkIntegrationGhcrDelete,

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
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					if value.(string) == "ghcr.io" {
						return nil, nil
					} else {
						return nil, []error{
							fmt.Errorf(
								"%s: can only be 'ghcr.io'", key,
							),
						}
					}
				},
			},
			"registry_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
						"ssl": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
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
			"limit_by_label": {
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
				Optional:    true,
				Description: "A list of key/value labels to limit the assessment of images",
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

func resourceLaceworkIntegrationGhcrCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GhcrContainerRegistry,
		api.GhcrData{
			LimitByTag:            castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:          castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:            castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:           d.Get("limit_num_imgs").(int),
			RegistryDomain:        d.Get("registry_domain").(string),
			RegistryNotifications: d.Get("registry_notifications").(bool),
			Credentials: api.GhcrCredentials{
				Username: d.Get("credentials.0.username").(string),
				Password: d.Get("credentials.0.password").(string),
				Ssl:      d.Get("credentials.0.ssl").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.GhcrContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.Create(data)
	if err != nil {
		return err
	}

	d.SetId(response.Data.IntgGuid)
	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)

	log.Printf("[INFO] Created ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), response.Data.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationGhcrRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetGhcr(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)

	creds := struct {
		username string
		password string
		ssl      bool
	}{username: response.Data.Data.Credentials.Username,
		password: response.Data.Data.Credentials.Password,
		ssl:      response.Data.Data.Credentials.Ssl}
	d.Set("credentials", creds)
	d.Set("registry_domain", response.Data.Data.RegistryDomain)
	d.Set("registry_notifications", response.Data.Data.RegistryNotifications)
	d.Set("limit_num_imgs", response.Data.Data.LimitNumImg)
	d.Set("limit_by_tags", response.Data.Data.LimitByTag)
	d.Set("limit_by_repositories", response.Data.Data.LimitByRep)
	d.Set("limit_by_label", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))

	log.Printf("[INFO] Read ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), response.Data.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationGhcrUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GhcrContainerRegistry,
		api.GhcrData{
			LimitByTag:            castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:          castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:            castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:           d.Get("limit_num_imgs").(int),
			RegistryDomain:        d.Get("registry_domain").(string),
			RegistryNotifications: d.Get("registry_notifications").(bool),
			Credentials: api.GhcrCredentials{
				Username: d.Get("credentials.0.username").(string),
				Password: d.Get("credentials.0.password").(string),
				Ssl:      d.Get("credentials.0.ssl").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.GhcrContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateGhcr(data)
	if err != nil {
		return err
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)

	log.Printf("[INFO] Updated ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationGhcrDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), d.Id())
	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted ContVulnCfg integration for %s registry type with guid %s\n",
		api.GhcrContainerRegistry.String(), d.Id())
	return nil
}
