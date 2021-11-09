package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationGcr() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGcrCreate,
		Read:   resourceLaceworkIntegrationGcrRead,
		Update: resourceLaceworkIntegrationGcrUpdate,
		Delete: resourceLaceworkIntegrationGcrDelete,

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
					switch value.(string) {
					case "gcr.io", "us.gcr.io", "eu.gcr.io", "asia.gcr.io":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be one of 'gcr.io', 'us.gcr.io', 'eu.gcr.io', or 'asia.gcr.io'", key,
							),
						}
					}
				},
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// @afiune we can't compare this element since our API, for security reasons,
								// does NOT return the private key configured in the Lacework server. So if
								// any other element changed from the credentials then we trigger a diff
								return !d.HasChanges(
									"name", "limit_by_tag", "limit_by_label", "org_level", "enabled",
									"credentials.0.client_id", "credentials.0.private_key_id",
									"credentials.0.client_email", "limit_by_repos", "limit_num_imgs",
									"limit_by_tags", "limit_by_labels", "limit_by_repositories",
								)
							},
						},
					},
				},
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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
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

func resourceLaceworkIntegrationGcrCreate(d *schema.ResourceData, meta interface{}) error {
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

	data := api.NewGcrRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:       limitByTags,
			LimitByLabel:     limitByLabels,
			LimitByRep:       limitByRepos,
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.ContainerRegCreds{
				ClientID:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), data)
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
			api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), integration.IntgGuid)
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationGcrRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), d.Id())
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

			creds := make(map[string]string)
			creds["client_id"] = integration.Data.Credentials.ClientID
			creds["client_email"] = integration.Data.Credentials.ClientEmail
			creds["private_key_id"] = integration.Data.Credentials.PrivateKeyID
			d.Set("credentials", []map[string]string{creds})
			d.Set("registry_domain", integration.Data.RegistryDomain)
			d.Set("limit_num_imgs", integration.Data.LimitNumImg)
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

			if _, ok := d.GetOk("limit_by_repositories"); ok {
				d.Set("limit_by_repositories", strings.Split(integration.Data.LimitByRep, ","))
			} else {
				d.Set("limit_by_repos", integration.Data.LimitByRep)
			}

			log.Printf("[INFO] Read %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGcrUpdate(d *schema.ResourceData, meta interface{}) error {
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

	data := api.NewGcrRegistryIntegration(d.Get("name").(string),
		api.ContainerRegData{
			LimitByTag:       limitByTags,
			LimitByLabel:     limitByLabels,
			LimitByRep:       limitByRepos,
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.ContainerRegCreds{
				ClientID:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration %s registry type with data:\n%+v\n",
		api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), data)
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
				api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), d.Id())

			return nil
		}
	}

	return nil
}

func resourceLaceworkIntegrationGcrDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), d.Id())

	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.GcrRegistry.String(), d.Id())

	return nil
}
