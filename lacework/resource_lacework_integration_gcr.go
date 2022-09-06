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

func resourceLaceworkIntegrationGcrCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GcpGarContainerRegistry,
		api.GcpGcrData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayOfKeyValueMap(d, "limit_by_labels"),
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

	log.Printf("[INFO] Creating %s registry type with data:\n%+v\n", api.GcpGcrContainerRegistry.String(), data)
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

	log.Printf("[INFO] Created %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), integration.IntgGuid)
	return nil

	return nil
}

func resourceLaceworkIntegrationGcrRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetGcpGcr(d.Id())

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

		creds := make(map[string]string)
		creds["client_id"] = integration.Data.Credentials.ClientID
		creds["client_email"] = integration.Data.Credentials.ClientEmail
		creds["private_key_id"] = integration.Data.Credentials.PrivateKeyID
		d.Set("credentials", []map[string]string{creds})
		d.Set("registry_domain", integration.Data.RegistryDomain)
		d.Set("limit_num_imgs", integration.Data.LimitNumImg)
		d.Set("non_os_package_support", integration.Data.NonOSPackageEval)

		d.Set("limit_by_tags", integration.Data.LimitByTag)
		d.Set("limit_by_repositories", integration.Data.LimitByRep)
		d.Set("limit_by_labels", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(integration.Data.LimitByLabel))

		log.Printf("[INFO] Read %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGcrUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.GcpGcrContainerRegistry,
		api.GcpGarData{
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayOfKeyValueMap(d, "limit_by_labels"),
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

	log.Printf("[INFO] Updating %s registry type with data:\n%+v\n", api.GcpGcrContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateGcpGcr(data)
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

		log.Printf("[INFO] Updated %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), d.Id())

		return nil
	}

	return nil
}

func resourceLaceworkIntegrationGcrDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), d.Id())

	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s registry type with guid: %v\n", api.GcpGcrContainerRegistry.String(), d.Id())

	return nil
}
