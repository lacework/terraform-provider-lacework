package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAzureActivityLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAzureActivityLogCreate,
		Read:   resourceLaceworkIntegrationAzureActivityLogRead,
		Update: resourceLaceworkIntegrationAzureActivityLogUpdate,
		Delete: resourceLaceworkIntegrationAzureActivityLogDelete,

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
			"queue_url": {
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
									"name", "tenant_id", "org_level", "queue_url",
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

func resourceLaceworkIntegrationAzureActivityLogCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		azure    = api.NewAzureIntegration(d.Get("name").(string),
			api.AzureActivityLogIntegration,
			api.AzureIntegrationData{
				TenantID: d.Get("tenant_id").(string),
				QueueUrl: d.Get("queue_url").(string),
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
	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AzureActivityLogIntegration.String(), azure)
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

	log.Printf("[INFO] Created %s integration with guid: %v\n",
		api.AzureActivityLogIntegration.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationAzureActivityLogRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.AzureActivityLogIntegration.String(), d.Id())
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
			d.Set("queue_url", integration.Data.QueueUrl)
			d.Set("tenant_id", integration.Data.TenantID)

			log.Printf("[INFO] Read %s integration with guid: %v\n", api.AzureActivityLogIntegration.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAzureActivityLogUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		azure    = api.NewAzureIntegration(d.Get("name").(string),
			api.AzureActivityLogIntegration,
			api.AzureIntegrationData{
				TenantID: d.Get("tenant_id").(string),
				QueueUrl: d.Get("queue_url").(string),
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

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AzureActivityLogIntegration.String(), azure)
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

	log.Printf("[INFO] Updated %sw integration with guid: %v\n", api.AzureActivityLogIntegration.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAzureActivityLogDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.AzureActivityLogIntegration.String(), d.Id())
	_, err := lacework.Integrations.DeleteAzure(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.AzureActivityLogIntegration.String(), d.Id())
	return nil
}
