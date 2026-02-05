package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkAzureDspm() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAzureDspmCreate,
		Read:   resourceLaceworkAzureDspmRead,
		Update: resourceLaceworkAzureDspmUpdate,
		Delete: resourceLaceworkAzureDspmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the FortiCNAPP DSPM integration.",
			},
			"integration_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the Azure tenant where the DSPM scanner is deployed.",
			},
			"storage_account_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the storage account where the DSPM scanner writes results and state.",
			},
			"blob_container_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the blob container where the DSPM scanner writes results and state.",
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The credentials used by Lacework platform to access the Azure SP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The clientID of the Azure SP used by Lacework platform.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The client secret of the Azure SP used by Lacework platform.",
						},
					},
				},
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The regions where the DSPM scanner is deployed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"server_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of attempts to create the external integration.",
			},
		},
	}
}

func resourceLaceworkAzureDspmCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	azureDspmData := api.AzureDspmData{
		TenantID:          d.Get("tenant_id").(string),
		StorageAccountUrl: d.Get("storage_account_url").(string),
		BlobContainerName: d.Get("blob_container_name").(string),
		Credentials: api.AzureDspmCredentials{
			ClientId:     d.Get("credentials.0.client_id").(string),
			ClientSecret: d.Get("credentials.0.client_secret").(string),
		},
		Regions: d.Get("regions").([]string),
	}

	azureDspm := api.NewCloudAccount(d.Get("name").(string),
		api.AzureDspmCloudAccount,
		azureDspmData,
	)

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AzureDspmCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAzureDspm(azureDspm)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s cloud account integration: %s",
						api.AzureDspmCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AzureDspmCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s cloud account integration (retrying %d more time(s))",
				api.AzureDspmCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		d.SetId(cloudAccount.IntgGuid)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("server_token", cloudAccount.ServerToken)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.AzureDspmCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkAzureDspmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	guid := d.Id()

	log.Printf("[INFO] Reading DSPM configuration %s", guid)
	resp, err := client.V2.CloudAccounts.GetAzureDspm(guid)
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.Set("regions", []string{})

	cloudAccount := resp.Data
	dspmData := cloudAccount.Data
	if cloudAccount.IntgGuid == d.Id() {
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("tenant_id", dspmData.TenantID)
		d.Set("storage_account_url", dspmData.StorageAccountUrl)
		d.Set("blob_container_name", dspmData.BlobContainerName)
		creds := make(map[string]string)
		creds["client_id"] = cloudAccount.Data.Credentials.ClientId
		creds["client_secret"] = cloudAccount.Data.Credentials.ClientSecret
		d.Set("credentials", []map[string]string{creds})
		d.Set("regions", dspmData.Regions)
		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AzureDspmCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAzureDspmUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	azureDspmData := api.AzureDspmData{
		TenantID:          d.Get("tenant_id").(string),
		StorageAccountUrl: d.Get("storage_account_url").(string),
		BlobContainerName: d.Get("blob_container_name").(string),
		Credentials: api.AzureDspmCredentials{
			ClientId:     d.Get("credentials.0.client_id").(string),
			ClientSecret: d.Get("credentials.0.client_secret").(string),
		},
		Regions: d.Get("regions").([]string),
	}

	azureDspm := api.NewCloudAccount(d.Get("name").(string),
		api.AzureDspmCloudAccount,
		azureDspmData,
	)

	azureDspm.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AzureDspmCloudAccount.String(), azureDspmData)
	_, err := lacework.V2.CloudAccounts.UpdateAzureDspm(azureDspm)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkAzureDspmDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	return nil
}
