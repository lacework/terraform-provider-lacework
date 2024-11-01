package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkIntegrationAzureActivityLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAzureActivityLogCreate,
		Read:   resourceLaceworkIntegrationAzureActivityLogRead,
		Update: resourceLaceworkIntegrationAzureActivityLogUpdate,
		Delete: resourceLaceworkIntegrationAzureActivityLogDelete,

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
								return !d.HasChanges(
									"name", "tenant_id", "org_level", "queue_url",
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

func resourceLaceworkIntegrationAzureActivityLogCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		azure    = api.NewCloudAccount(d.Get("name").(string),
			api.AzureAlSeqCloudAccount,
			api.AzureAlSeqData{
				TenantID: d.Get("tenant_id").(string),
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AzureAlSeqCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientSecret: d.Get("credentials.0.client_secret").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		azure.Enabled = 0
	}

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AzureAlSeqCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(azure)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("error creating %s integration: %s",
						api.AzureAlSeqCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AzureAlSeqCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"unable to create %s integration (retrying %d more time(s))",
				api.AzureAlSeqCloudAccount.String(), retries,
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
			api.AzureAlSeqCloudAccount.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAzureActivityLogRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.AzureAlSeqCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAzureAlSeq(d.Id())
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

		creds := make(map[string]string)
		creds["client_id"] = integration.Data.Credentials.ClientID
		d.Set("credentials", []map[string]string{creds})
		d.Set("queue_url", integration.Data.QueueUrl)
		d.Set("tenant_id", integration.Data.TenantID)

		log.Printf("[INFO] Read %s integration with guid: %v\n", api.AzureAlSeqCloudAccount.String(), integration.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAzureActivityLogUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		azure    = api.NewCloudAccount(d.Get("name").(string),
			api.AzureAlSeqCloudAccount,
			api.AzureAlSeqData{
				TenantID: d.Get("tenant_id").(string),
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AzureAlSeqCredentials{
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

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AzureAlSeqCloudAccount.String(), azure)
	response, err := lacework.V2.CloudAccounts.UpdateAzureAlSeq(azure)
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

	log.Printf("[INFO] Updated %sw integration with guid: %v\n", api.AzureAlSeqCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAzureActivityLogDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.AzureAlSeqCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.AzureAlSeqCloudAccount.String(), d.Id())
	return nil
}
