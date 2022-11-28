package lacework

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
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
				Type:     schema.TypeString,
				Required: true,
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
							Type:     schema.TypeString,
							Required: true,
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
									"name", "resource_level", "resource_id", "org_level", "enabled",
									"credentials.0.client_id",
									"credentials.0.client_email",
								)
							},
						},
						"token_uri": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "https://oauth2.googleapis.com/token",
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
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scanning_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scan_frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  24,
			},
			"scan_containers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"scan_host_vulnerabilities": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"query_text": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"filter_list": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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
			FilterList:              d.Get("filter_list").(string),
		},
	)

	f, _ := os.Create("/tmp/yourfile")
	defer f.Close()

	w := bufio.NewWriter(f)
	//choose random number for recipe

	_, _ = fmt.Fprintf(w, "%v\n", data)
	fmt.Printf("Data is %+v", data)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.GcpSidekickCloudAccount.String())
		log.Printf("[INFO] Creating %v integration\n", data)
		response, err := lacework.V2.CloudAccounts.Create(data)

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
		d.Set("filter_list", integration.Data.FilterList)

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
			FilterList:              d.Get("filter_list").(string),
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
