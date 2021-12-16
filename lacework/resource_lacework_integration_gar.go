package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

var GarDomainRegistries = []string{
	"northamerica-northeast1-docker.pkg.dev",
	"us-central1-docker.pkg.dev",
	"us-east1-docker.pkg.dev",
	"us-east4-docker.pkg.dev",
	"us-west1-docker.pkg.dev",
	"us-west2-docker.pkg.dev",
	"us-west3-docker.pkg.dev",
	"us-west4-docker.pkg.dev",
	"southamerica-east1-docker.pkg.dev",
	"europe-north1-docker.pkg.dev",
	"europe-west1-docker.pkg.dev",
	"europe-west2-docker.pkg.dev",
	"europe-west3-docker.pkg.dev",
	"europe-west4-docker.pkg.dev",
	"europe-west6-docker.pkg.dev",
	"asia-east1-docker.pkg.dev",
	"asia-east2-docker.pkg.dev",
	"asia-northeast1-docker.pkg.dev",
	"asia-northeast2-docker.pkg.dev",
	"asia-northeast3-docker.pkg.dev",
	"asia-south1-docker.pkg.dev",
	"asia-southeast1-docker.pkg.dev",
	"asia-southeast2-docker.pkg.dev",
	"australia-southeast1-docker.pkg.dev",
	"asia-docker.pkg.dev",
	"europe-docker.pkg.dev",
	"us-docker.pkg.dev",
}

func resourceLaceworkIntegrationGar() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGarCreate,
		Read:   resourceLaceworkIntegrationGarRead,
		Update: resourceLaceworkIntegrationGarUpdate,
		Delete: resourceLaceworkIntegrationGarDelete,

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
					if ContainsStr(GarDomainRegistries, value.(string)) {
						return nil, nil
					}

					return nil, []error{
						fmt.Errorf(
							"%s: can only be one of:\n%s", key, strings.Join(GarDomainRegistries, "\n"),
						),
					}
				},
			},
			"non_os_package_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable program language scanning",
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
						"client_email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key_id": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// @afiune we can't compare this element since our API, for security reasons,
								// does NOT return the private key configured in the Lacework server. So if
								// any other element changed from the credentials then we trigger a diff
								return !d.HasChanges(
									"name", "registry_domain", "enabled", "non_os_package_support",
									"limit_by_tags", "limit_by_label", "limit_by_repositories",
									"limit_num_imgs", "credentials.0.client_id", "credentials.0.client_email",
								)
							},
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
									"name", "registry_domain", "enabled", "non_os_package_support",
									"limit_by_tags", "limit_by_label", "limit_by_repositories",
									"limit_num_imgs", "credentials.0.client_id", "credentials.0.client_email",
								)
							},
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
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The maximum number of newest container images to assess per repository.",
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

func resourceLaceworkIntegrationGarCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GcpGarContainerRegistry,
		api.GcpGarData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.GcpCredentialsV2{
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

	log.Printf("[INFO] Creating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.GcpGarContainerRegistry.String(), data)
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
		api.GcpGarContainerRegistry.String(), response.Data.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationGarRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading ContVulnCfg integration for %s registry type with guid %s\n",
		api.GcpGarContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetGcpGar(d.Id())
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
	d.Set("non_os_package_support", response.Data.Data.NonOSPackageEval)

	creds := make(map[string]string)
	creds["client_id"] = response.Data.Data.Credentials.ClientID
	creds["client_email"] = response.Data.Data.Credentials.ClientEmail
	d.Set("credentials", []map[string]string{creds})
	d.Set("registry_domain", response.Data.Data.RegistryDomain)
	d.Set("limit_num_imgs", response.Data.Data.LimitNumImg)
	d.Set("limit_by_tags", response.Data.Data.LimitByTag)
	d.Set("limit_by_repositories", response.Data.Data.LimitByRep)
	d.Set("limit_by_label", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))

	log.Printf("[INFO] Read ContVulnCfg integration for %s registry type with guid %s\n",
		api.GcpGarContainerRegistry.String(), response.Data.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationGarUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GcpGarContainerRegistry,
		api.GcpGarData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
			Credentials: api.GcpCredentialsV2{
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

	log.Printf("[INFO] Updating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.GcpGarContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateGcpGar(data)
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
		api.GcpGarContainerRegistry.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationGarDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting ContVulnCfg integration for %s registry type with guid %s\n",
		api.GcpGarContainerRegistry.String(), d.Id())
	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted ContVulnCfg integration for %s registry type with guid %s\n",
		api.GcpGarContainerRegistry.String(), d.Id())
	return nil
}
