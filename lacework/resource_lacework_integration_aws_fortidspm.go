package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

// lacework_integration_aws_fortidspm registers a FortiDSPM-managed AWS DSPM
// cloud account. Unlike lacework_integration_aws_dspm there is no customer
// bucket and no cross-account role: FortiDSPM pushes the results to Lacework
// itself. Creating the account makes api-server ask FortiDSPM for a
// deployment, and the per-region activation tokens and AMIs it returns are
// exposed here so the scan engine module can boot the instances.
func resourceLaceworkIntegrationAwsFortiDspm() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsFortiDspmCreate,
		Read:   resourceLaceworkIntegrationAwsFortiDspmRead,
		Update: resourceLaceworkIntegrationAwsFortiDspmUpdate,
		Delete: resourceLaceworkIntegrationAwsFortiDspmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the FortiCNAPP DSPM integration.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the AWS account where the scan engines are deployed.",
			},
			"regions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "The regions where a scan engine is deployed.",
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
			"activation_tokens": {
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Description: "Single-use scan engine activation token per region, issued once on create.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"image_ids": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Scan engine AMI id per region.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceLaceworkIntegrationAwsFortiDspmCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	account := api.NewCloudAccount(d.Get("name").(string),
		api.AwsDspmCloudAccount,
		api.AwsFortiDspmData{
			AccountID: d.Get("account_id").(string),
			Regions:   castAttributeToStringSlice(d, "regions"),
		},
	)

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating FortiDSPM-managed %s cloud account integration\n", api.AwsDspmCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAwsFortiDspm(account)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s cloud account integration: %s",
						api.AwsDspmCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AwsDspmCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s cloud account integration (retrying %d more time(s))",
				api.AwsDspmCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		deployment := cloudAccount.FortiDspmDeployment
		if deployment == nil {
			// Nothing to retry: the account exists but api-server did not run
			// the FortiDSPM step, which means it is not configured for it.
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
		d.Set("activation_tokens", deployment.ActivationTokens())
		d.Set("image_ids", deployment.Images())

		log.Printf("[INFO] Created FortiDSPM-managed %s cloud account integration with guid: %v (deployment %s)\n",
			api.AwsDspmCloudAccount.String(), cloudAccount.IntgGuid, deployment.DeploymentID)
		return nil
	})
}

// Read confirms the integration still exists and refreshes what the API
// stores. The deployment attributes are issued once on create and are not
// returned by reads, so they are left as they are in state.
func resourceLaceworkIntegrationAwsFortiDspmRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsFortiDspm(d.Id())
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
	d.Set("account_id", cloudAccount.Data.AccountID)
	d.Set("regions", cloudAccount.Data.Regions)
	return nil
}

func resourceLaceworkIntegrationAwsFortiDspmUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	account := api.NewCloudAccount(d.Get("name").(string),
		api.AwsDspmCloudAccount,
		api.AwsFortiDspmData{
			AccountID: d.Get("account_id").(string),
			Regions:   castAttributeToStringSlice(d, "regions"),
		},
	)
	account.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	if _, err := lacework.V2.CloudAccounts.UpdateAwsDspm(account); err != nil {
		return err
	}
	return nil
}

func resourceLaceworkIntegrationAwsFortiDspmDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	if err := lacework.V2.CloudAccounts.Delete(d.Id()); err != nil {
		return err
	}
	return nil
}
