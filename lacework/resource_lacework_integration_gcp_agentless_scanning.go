package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationGcpAgentlessScanning() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGcpAgentlessScanningCreate,
		Read:   resourceLaceworkIntegrationGcpAgentlessScanningRead,
		Update: resourceLaceworkIntegrationGcpAgentlessScanningUpdate,
		Delete: resourceLaceworkIntegrationGcpAgentlessScanningDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkCloudAccount,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The integration name.",
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration.",
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of attempts to create the external integration.",
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client Id from credentials file.",
						},
						"private_key_id": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return !d.HasChanges(
									"name", "resource_level", "resource_id", "org_level", "enabled",
									"credentials.0.client_id",
									"credentials.0.client_email",
								)
							},
							Description: "Private Key Id from credentials file.",
						},
						"client_email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client email from credentials file.",
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
									"name", "resource_level", "resource_id", "org_level", "enabled",
									"credentials.0.client_id",
									"credentials.0.client_email",
								)
							},
							Description: "Private Key from credentials file.",
						},
						"token_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "https://oauth2.googleapis.com/token",
							Description: "Token URI from credentials file.",
						},
					},
				},
			},
			"resource_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.GcpProjectIntegration.String(),
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch strings.ToUpper(value.(string)) {
					case api.GcpProjectIntegration.String(),
						api.GcpOrganizationIntegration.String():
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf("%s: can only be either '%s' or '%s'",
								key,
								api.GcpProjectIntegration.String(),
								api.GcpOrganizationIntegration.String()),
						}
					}
				},
				Description: "Integration level - ORGANIZATION / PROJECT.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization Id or Project Id.",
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
			"server_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bucket containing analysis results shared with Lacework platform.",
			},
			"scanning_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID where scanner is deployed.",
			},
			"scan_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     24,
				Description: "How often in hours the scan will run in hours.",
			},
			"scan_containers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to includes scanning for containers.",
			},
			"scan_host_vulnerabilities": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to includes scanning for host vulnerabilities.",
			},
			"query_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The LQL query text.",
			},
			"filter_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Default:     nil,
				Description: "List of Projects to specifically include/exclude.",
			},
		},
	}
}

func resourceLaceworkIntegrationGcpAgentlessScanningCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		retries       = d.Get("retries").(int)
		resourceLevel = api.GcpProjectIntegration
	)

	if strings.ToUpper(
		d.Get("resource_level").(string),
	) == api.GcpOrganizationIntegration.String() {
		resourceLevel = api.GcpOrganizationIntegration
	}
	log.Printf("[INFO] Creating %s integration\n", api.GcpSidekickCloudAccount.String())

	data := api.NewCloudAccount(d.Get("name").(string),
		api.GcpSidekickCloudAccount,
		api.GcpSidekickData{
			ID:     d.Get("resource_id").(string),
			IDType: resourceLevel.String(),
			Credentials: api.GcpSidekickCredentials{
				ClientID:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
				TokenUri:     d.Get("credentials.0.token_uri").(string),
			},
			SharedBucket:            d.Get("bucket_name").(string),
			ScanningProjectId:       d.Get("scanning_project_id").(string),
			ScanFrequency:           d.Get("scan_frequency").(int),
			ScanContainers:          d.Get("scan_containers").(bool),
			ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
			QueryText:               d.Get("query_text").(string),
			FilterList:              strings.Join(castAttributeToStringSlice(d, "filter_list"), ", "),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.GcpSidekickCloudAccount.String())
		log.Printf("[INFO] Creating %v integration\n", data)
		response, err := lacework.V2.CloudAccounts.CreateGcpSidekick(data)

		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.GcpSidekickCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.GcpSidekickCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.GcpSidekickCloudAccount.String(), retries,
			))
		}

		integration := response.Data
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("resource_id", integration.IntgGuid)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)
		d.Set("server_token", integration.ServerToken)
		d.Set("uri", integration.Uri)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.GcpSidekickCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationGcpAgentlessScanningRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.GcpSidekickCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetGcpSidekick(d.Id())
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
		d.Set("server_token", integration.ServerToken)
		d.Set("uri", integration.Uri)

		creds := make(map[string]string)
		creds["client_id"] = integration.Data.Credentials.ClientID
		creds["client_email"] = integration.Data.Credentials.ClientEmail
		creds["private_key_id"] = integration.Data.Credentials.PrivateKeyID
		creds["token_uri"] = integration.Data.Credentials.TokenUri
		d.Set("credentials", []map[string]string{creds})
		d.Set("resource_level", integration.Data.IDType)
		d.Set("resource_id", integration.Data.ID)
		d.Set("bucket_name", integration.Data.SharedBucket)
		d.Set("scanning_project_id", integration.Data.ScanningProjectId)
		d.Set("scan_frequency", integration.Data.ScanFrequency)
		d.Set("scan_containers", integration.Data.ScanContainers)
		d.Set("scan_host_vulnerabilities", integration.Data.ScanHostVulnerabilities)
		d.Set("query_text", integration.Data.QueryText)

		filter_list := strings.Split(integration.Data.FilterList, ",")
		if integration.Data.FilterList != "" && len(filter_list) > 0 {
			var trimmed_filter_list []string
			for _, elem := range filter_list {
				trimmed_filter_list = append(trimmed_filter_list, strings.TrimSpace(elem))
			}
			d.Set("filter_list", trimmed_filter_list)
		}

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.GcpSidekickCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGcpAgentlessScanningUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		resourceLevel = api.GcpProjectIntegration
	)

	if strings.ToUpper(d.Get("resource_level").(string)) == api.GcpOrganizationIntegration.String() {
		resourceLevel = api.GcpOrganizationIntegration
	}

	data := api.NewCloudAccount(d.Get("name").(string),
		api.GcpSidekickCloudAccount,
		api.GcpSidekickData{
			ID:     d.Get("resource_id").(string),
			IDType: resourceLevel.String(),
			Credentials: api.GcpSidekickCredentials{
				ClientID:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
				TokenUri:     d.Get("credentials.0.token_uri").(string),
			},
			SharedBucket:            d.Get("bucket_name").(string),
			ScanningProjectId:       d.Get("scanning_project_id").(string),
			ScanFrequency:           d.Get("scan_frequency").(int),
			ScanContainers:          d.Get("scan_containers").(bool),
			ScanHostVulnerabilities: d.Get("scan_host_vulnerabilities").(bool),
			QueryText:               d.Get("query_text").(string),
			FilterList:              strings.Join(castAttributeToStringSlice(d, "filter_list"), ", "),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.GcpSidekickCloudAccount.String(), data)
	response, err := lacework.V2.CloudAccounts.UpdateGcpSidekick(data)
	if err != nil {
		return err
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("resource_level", integration.Data.IDType)
	d.Set("resource_id", integration.Data.ID)

	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)
	d.Set("server_token", integration.ServerToken)
	d.Set("uri", integration.Uri)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.GcpSidekickCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationGcpAgentlessScanningDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.GcpSidekickCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.GcpSidekickCloudAccount.String(), d.Id())
	return nil
}
