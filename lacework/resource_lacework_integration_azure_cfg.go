package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
								if d.HasChanges(
									"name", "tenant_id", "org_level",
									"enabled", "credentials.0.client_id",
								) {
									return false
								}
								return true
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
		azure    = api.NewAzureIntegration(d.Get("name").(string),
			api.AzureCfgIntegration,
			api.AzureIntegrationData{
				TenantID: d.Get("tenant_id").(string),
				Credentials: api.AzureIntegrationCreds{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientSecret: d.Get("credentials.0.client_secret").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		azure.Enabled = 0
	}

	// @afiune should we do this if there is sensitive information?
	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AzureCfgIntegration.String(), azure)
	response, err := lacework.Integrations.CreateAzure(azure)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAzureIntegrationResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.AzureCfgIntegration.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationAzureCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.AzureCfgIntegration.String(), d.Id())
	response, err := lacework.Integrations.GetAzure(d.Id())
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
			d.Set("credentials", []map[string]string{creds})
			d.Set("tenant_id", integration.Data.TenantID)

			log.Printf("[INFO] Read %s integration with guid: %v\n", api.AzureCfgIntegration.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAzureCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		azure    = api.NewAzureIntegration(d.Get("name").(string),
			api.AzureCfgIntegration,
			api.AzureIntegrationData{
				TenantID: d.Get("tenant_id").(string),
				Credentials: api.AzureIntegrationCreds{
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

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AzureCfgIntegration.String(), azure)
	response, err := lacework.Integrations.UpdateAzure(azure)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAzureIntegrationResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point of time, we know the data field has a single value
	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.AzureCfgIntegration.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAzureCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.AzureCfgIntegration.String(), d.Id())
	_, err := lacework.Integrations.DeleteAzure(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.AzureCfgIntegration.String(), d.Id())
	return nil
}

// validateAzureIntegrationResponse checks weather or not the server response has
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
