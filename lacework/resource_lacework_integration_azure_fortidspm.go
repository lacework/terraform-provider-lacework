package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

// lacework_integration_azure_fortidspm registers a FortiDSPM-managed Azure
// DSPM cloud account. FortiDSPM returns, per region, a single-use activation
// token and a URL to the scan engine image VHD; the scan engine module copies
// the VHD into the customer's subscription and boots from it.
func resourceLaceworkIntegrationAzureFortiDspm() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAzureFortiDspmCreate,
		Read:   resourceLaceworkIntegrationAzureFortiDspmRead,
		Update: resourceLaceworkIntegrationAzureFortiDspmUpdate,
		Delete: resourceLaceworkIntegrationAzureFortiDspmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the FortiCNAPP DSPM integration.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Azure tenant where the scan engines are deployed.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The Azure subscription where the scan engines are deployed. Empty for a tenant-level integration.",
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "The Azure locations where a scan engine is deployed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of attempts to create the integration.",
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FortiDSPM deployment id.",
			},
			"deployment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"env_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FortiDSPM environment id.",
			},
			"token_expires_in": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Seconds the activation tokens stay valid after creation.",
			},
			"image_url_expires_in": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Seconds the image URLs stay valid after creation.",
			},
			"activation_tokens": {
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Description: "Single-use scan engine activation token per region, issued once on create.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"image_urls": {
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Description: "Signed URL of the scan engine image VHD per region, issued once on create.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"hyperv_generations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Hyper-V generation (V1 or V2) of the scan engine image per region.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func azureFortiDspmData(d *schema.ResourceData) api.AzureFortiDspmData {
	return api.AzureFortiDspmData{
		TenantID:       d.Get("tenant_id").(string),
		SubscriptionID: d.Get("subscription_id").(string),
		Regions:        castAttributeToStringSlice(d, "regions"),
	}
}

func resourceLaceworkIntegrationAzureFortiDspmCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	account := api.NewCloudAccount(d.Get("name").(string), api.AzureDspmCloudAccount, azureFortiDspmData(d))

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating FortiDSPM-managed %s cloud account integration\n", api.AzureDspmCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAzureFortiDspm(account)
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
		deployment := cloudAccount.FortiDspmDeployment
		if deployment == nil {
			if delErr := lacework.V2.CloudAccounts.Delete(cloudAccount.IntgGuid); delErr != nil {
				log.Printf("[WARN] Failed to delete integration %s after missing FortiDSPM deployment: %s",
					cloudAccount.IntgGuid, delErr)
			}
			return retry.NonRetryableError(fmt.Errorf(
				"the API created integration %s but returned no FortiDSPM deployment; "+
					"FortiDSPM is not enabled for this account", cloudAccount.IntgGuid))
		}

		d.SetId(cloudAccount.IntgGuid)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("deployment_id", deployment.DeploymentID)
		d.Set("deployment_name", deployment.DeploymentName)
		d.Set("env_id", deployment.EnvID)
		d.Set("token_expires_in", deployment.TokenExpiresIn)
		d.Set("image_url_expires_in", deployment.ImageURLExpiresIn)
		d.Set("activation_tokens", deployment.ActivationTokens())
		d.Set("image_urls", deployment.Images())
		d.Set("hyperv_generations", deployment.HypervGenerations())

		log.Printf("[INFO] Created FortiDSPM-managed %s cloud account integration with guid: %v (deployment %s)\n",
			api.AzureDspmCloudAccount.String(), cloudAccount.IntgGuid, deployment.DeploymentID)
		return nil
	})
}

func resourceLaceworkIntegrationAzureFortiDspmRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAzureFortiDspm(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	cloudAccount := response.Data
	if cloudAccount.IntgGuid != d.Id() {
		d.SetId("")
		return nil
	}
	d.Set("intg_guid", cloudAccount.IntgGuid)
	d.Set("name", cloudAccount.Name)
	d.Set("tenant_id", cloudAccount.Data.TenantID)
	d.Set("subscription_id", cloudAccount.Data.SubscriptionID)
	d.Set("regions", cloudAccount.Data.Regions)
	return nil
}

func resourceLaceworkIntegrationAzureFortiDspmUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	account := api.NewCloudAccount(d.Get("name").(string), api.AzureDspmCloudAccount, azureFortiDspmData(d))
	account.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	if _, err := lacework.V2.CloudAccounts.UpdateAzureDspm(account); err != nil {
		return err
	}
	return nil
}

func resourceLaceworkIntegrationAzureFortiDspmDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AzureDspmCloudAccount.String(), d.Id())
	if err := lacework.V2.CloudAccounts.Delete(d.Id()); err != nil {
		return err
	}
	return nil
}
