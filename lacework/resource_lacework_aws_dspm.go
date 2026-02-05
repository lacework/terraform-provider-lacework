package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkAwsDspm() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAwsDspmCreate,
		Read:   resourceLaceworkAwsDspmRead,
		Update: resourceLaceworkAwsDspmUpdate,
		Delete: resourceLaceworkAwsDspmDelete,
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
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the AWS account where the DSPM scanner is deployed.",
			},
			"storage_bucket_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the s3 bucket where the DSPM scanner writes results and state.",
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The credentials needed by the integration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The external id.",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ARN of the IAM role used by the DSPM .",
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

func resourceLaceworkAwsDspmCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
	)

	awsDspmData := api.AwsDspmData{
		AccountID: d.Get("account_id").(string),
		BucketArn: d.Get("storage_bucket_arn").(string),
		CrossAccountCreds: api.AwsDspmCrossAccountCredentials{
			ExternalID: d.Get("credentials.0.external_id").(string),
			RoleArn:    d.Get("credentials.0.role_arn").(string),
		},
		Regions: d.Get("regions").([]string),
	}

	awsDspm := api.NewCloudAccount(d.Get("name").(string),
		api.AwsDspmCloudAccount,
		awsDspmData,
	)

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AwsDspmCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.CreateAwsDspm(awsDspm)
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
		d.SetId(cloudAccount.IntgGuid)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("server_token", cloudAccount.ServerToken)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.AwsDspmCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkAwsDspmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	guid := d.Id()

	log.Printf("[INFO] Reading DSPM configuration %s", guid)
	resp, err := client.V2.CloudAccounts.GetAwsDspm(guid)
	if err != nil {
		return resourceNotFound(d, err)
	}

	cloudAccount := resp.Data
	dspmData := cloudAccount.Data
	if cloudAccount.IntgGuid == d.Id() {
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("account_id", dspmData.AccountID)
		d.Set("storage_bucket_arn", dspmData.BucketArn)
		creds := make(map[string]string)
		creds["role_arn"] = dspmData.CrossAccountCreds.RoleArn
		creds["external_id"] = dspmData.CrossAccountCreds.ExternalID
		d.Set("credentials", []map[string]string{creds})
		d.Set("regions", dspmData.Regions)

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AwsDspmCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAwsDspmUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	awsDspmData := api.AwsDspmData{
		AccountID: d.Get("account_id").(string),
		BucketArn: d.Get("storage_bucket_arn").(string),
		CrossAccountCreds: api.AwsDspmCrossAccountCredentials{
			ExternalID: d.Get("credentials.0.external_id").(string),
			RoleArn:    d.Get("credentials.0.role_arn").(string),
		},
		Regions: d.Get("regions").([]string),
	}

	awsDspm := api.NewCloudAccount(d.Get("name").(string),
		api.AwsDspmCloudAccount,
		awsDspmData,
	)

	awsDspm.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsDspmCloudAccount.String(), awsDspmData)
	_, err := lacework.V2.CloudAccounts.UpdateAwsDspm(awsDspm)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkAwsDspmDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AwsDspmCloudAccount.String(), d.Id())
	return nil
}
