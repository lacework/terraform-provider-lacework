package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAzureCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAzureCfgCreate,
		Read:   resourceLaceworkIntegrationAzureCfgRead,
		Update: resourceLaceworkIntegrationAzureCfgUpdate,
		Delete: resourceLaceworkIntegrationAzureCfgDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"client_secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// @afiune we can't compare this element since our API, for security reasons,
								// does NOT return the client secret configured in the Lacework server. So if
								// any other element changed from the credentials then we trigger a diff
								return !d.HasChanges(
									"name", "tenant_id", "org_level",
									"enabled", "credentials.0.client_id",
								)
							},
						},
					},
				},
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

func resourceLaceworkIntegrationAzureCfgCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		azure    = api.NewCloudAccount(d.Get("name").(string),
			api.AwsCfgCloudAccount,
			api.AzureCfgData{
				TenantID: d.Get("tenant_id").(string),
				Credentials: api.AzureCfgCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientSecret: d.Get("credentials.0.client_secret").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		azure.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AwsCfgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(azure)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("error creating %s integration: %s",
						api.AwsCfgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AwsCfgCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"unable to create %s integration (retrying %d more time(s))",
				api.AwsCfgCloudAccount.String(), retries,
			))
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

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.AwsCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAzureCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAzureCfg(d.Id())
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
		d.Set("credentials", []map[string]string{creds})
		d.Set("tenant_id", integration.Data.TenantID)

		log.Printf("[INFO] Read %s integration with guid: %v\n", api.AwsCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAzureCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		azure    = api.NewCloudAccount(d.Get("name").(string),
			api.AwsCfgCloudAccount,
			api.AzureCfgData{
				TenantID: d.Get("tenant_id").(string),
				Credentials: api.AzureCfgCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientSecret: d.Get("credentials.0.client_secret").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		azure.Enabled = 0
	}

	azure.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AwsCfgCloudAccount.String(), azure)
	response, err := lacework.V2.CloudAccounts.UpdateAzureCfg(azure)
	if err != nil {
		return err
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAzureCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AwsCfgCloudAccount.String(), d.Id())
	return nil
}

// validateAzureIntegrationResponse checks whether or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateAzureIntegrationResponse(response *api.AzureIntegrationsResponse) error {
	if len(response.Data) == 0 {
		// @afiune this edge case should never happen, if we land here it means that
		// something went wrong in the server side of things (Lacework API), so let
		// us inform that to our users
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		// @afiune if we are creating a single integration and the server returns
		// more than one integration inside the 'data' field, it is definitely another
		// edge case that should never happen
		msg := `
There is more that one integration inside the server response data.

List of integrations:
`
		for _, integration := range response.Data {
			msg = msg + fmt.Sprintf("\t%s: %s\n", integration.IntgGuid, integration.Name)
		}
		msg = msg + unexpectedBehaviorMsg()
		return fmt.Errorf(msg)
	}

	return nil
}
