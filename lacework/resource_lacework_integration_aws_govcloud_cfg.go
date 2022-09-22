package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsGovCloudCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsGovCloudCfgCreate,
		Read:   resourceLaceworkIntegrationAwsGovCloudCfgRead,
		Update: resourceLaceworkIntegrationAwsGovCloudCfgUpdate,
		Delete: resourceLaceworkIntegrationAwsGovCloudCfgDelete,

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
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS Account ID",
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AWS access key ID",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The AWS secret key for the specified AWS access key",
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

func resourceLaceworkIntegrationAwsGovCloudCfgCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		aws      = api.NewCloudAccount(d.Get("name").(string),
			api.AwsUsGovCfgCloudAccount,
			api.AwsUsGovCfgData{
				Credentials: api.AwsUsGovCfgCredentials{
					AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
					SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
					AwsAccountID:    d.Get("account_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AwsUsGovCfgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(aws)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("error creating %s integration: %s",
						api.AwsUsGovCfgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AwsUsGovCfgCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"unable to create %s integration (retrying %d more time(s))",
				api.AwsUsGovCfgCloudAccount.String(), retries,
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
			api.AwsUsGovCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsGovCloudCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AwsUsGovCfgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsUsGovCfg(d.Id())
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
		d.Set("account_id", integration.Data.Credentials.AwsAccountID)

		creds := make(map[string]string)
		credentials := integration.Data.Credentials
		creds["access_key_id"] = credentials.AccessKeyID
		creds["secret_access_key"] = credentials.SecretAccessKey
		d.Set("credentials", []map[string]string{creds})

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.AwsUsGovCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsGovCloudCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewCloudAccount(d.Get("name").(string),
			api.AwsCfgCloudAccount,
			api.AwsUsGovCfgData{
				Credentials: api.AwsUsGovCfgCredentials{
					AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
					SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
					AwsAccountID:    d.Get("account_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AwsUsGovCfgCloudAccount.String(), aws)
	response, err := lacework.V2.CloudAccounts.UpdateAwsUsGovCfg(aws)
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
		api.AwsUsGovCfgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsGovCloudCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AwsUsGovCfgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AwsUsGovCfgCloudAccount.String(), d.Id())
	return nil
}

// validateAwsIntegrationResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateAwsIntegrationResponse(response *api.AwsIntegrationsResponse) error {
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

func unexpectedBehaviorMsg() string {
	return `
This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
}
