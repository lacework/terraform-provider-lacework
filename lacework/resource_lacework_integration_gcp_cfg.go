package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkIntegrationGcpCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGcpCfgCreate,
		Read:   resourceLaceworkIntegrationGcpCfgRead,
		Update: resourceLaceworkIntegrationGcpCfgUpdate,
		Delete: resourceLaceworkIntegrationGcpCfgDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
		},
	}
}

func resourceLaceworkIntegrationGcpCfgCreate(d *schema.ResourceData, meta interface{}) error {
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

	data := api.NewCloudAccount(d.Get("name").(string),
		api.GcpCfgCloudAccount,
		api.GcpCfgData{
			ID:     d.Get("resource_id").(string),
			IDType: resourceLevel.String(),
			Credentials: api.GcpCfgCredentials{
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

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.GcpCfgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(data)

		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.GcpCfgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.GcpCfgCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.GcpCfgCloudAccount.String(), retries,
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
			api.GcpCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationGcpCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.GcpCfgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetGcpCfg(d.Id())
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
		d.Set("credentials", []map[string]string{creds})
		d.Set("resource_level", integration.Data.IDType)
		d.Set("resource_id", integration.Data.ID)

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.GcpCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGcpCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		resourceLevel = api.GcpProjectIntegration
	)

	if strings.ToUpper(d.Get("resource_level").(string)) == api.GcpOrganizationIntegration.String() {
		resourceLevel = api.GcpOrganizationIntegration
	}

	data := api.NewCloudAccount(d.Get("name").(string),
		api.GcpCfgCloudAccount,
		api.GcpCfgData{
			ID:     d.Get("resource_id").(string),
			IDType: resourceLevel.String(),
			Credentials: api.GcpCfgCredentials{
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

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.GcpCfgCloudAccount.String(), data)
	response, err := lacework.V2.CloudAccounts.UpdateGcpCfg(data)
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
		api.GcpCfgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationGcpCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.GcpCfgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.GcpCfgCloudAccount.String(), d.Id())
	return nil
}
