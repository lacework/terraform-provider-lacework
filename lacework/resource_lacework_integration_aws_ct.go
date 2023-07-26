package lacework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsCloudTrail() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationAwsCloudTrailCreate,
		Read:     resourceLaceworkIntegrationAwsCloudTrailRead,
		Update:   resourceLaceworkIntegrationAwsCloudTrailUpdate,
		Delete:   resourceLaceworkIntegrationAwsCloudTrailDelete,
		Schema:   awsCloudTrailIntegrationSchema,
		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

var awsCloudTrailIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The integration name.",
	},
	"intg_guid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "The state of the external integration.",
	},
	"retries": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     5,
		Description: "The number of attempts to create the external integration.",
	},
	"queue_url": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The SQS Queue URL.",
	},
	"credentials": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Required:    true,
		Description: "The credentials needed by the integration.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:     schema.TypeString,
					Required: true,
				},
				"external_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	"org_account_mappings": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Mapping of AWS accounts to Lacework accounts within a Lacework organization.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_lacework_account": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The default Lacework account name where any non-mapped AWS account will appear",
				},
				"mapping": {
					Type:        schema.TypeSet,
					Required:    true,
					Description: "A map of AWS accounts to Lacework account. This can be specified multiple times to map multiple Lacework accounts.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"lacework_account": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Lacework account name where the CloudTrail activity from the selected AWS accounts will appear.",
							},
							"aws_accounts": {
								Type:        schema.TypeSet,
								Elem:        &schema.Schema{Type: schema.TypeString},
								MinItems:    1,
								Required:    true,
								Description: "The list of AWS account IDs to map.",
							},
						},
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
}

func resourceLaceworkIntegrationAwsCloudTrailCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework     = meta.(*api.Client)
		retries      = d.Get("retries").(int)
		awsCtSqsData = api.AwsCtSqsData{
			QueueUrl: d.Get("queue_url").(string),
			Credentials: api.AwsCtSqsCredentials{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
		}
	)
	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d)
	if !accountMapFile.Empty() {
		accountMapFileBytes, err := json.Marshal(accountMapFile)
		if err != nil {
			return err
		}

		awsCtSqsData.EncodeAccountMappingFile(accountMapFileBytes)
	}

	awsCtSqs := api.NewCloudAccount(d.Get("name").(string),
		api.AwsCtSqsCloudAccount,
		awsCtSqsData,
	)

	if !d.Get("enabled").(bool) {
		awsCtSqs.Enabled = 0
	}

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.AwsCtSqsCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(awsCtSqs)
		if err != nil {
			if retries <= 0 {
				return retry.NonRetryableError(
					fmt.Errorf("Error creating %s cloud account integration: %s",
						api.AwsCtSqsCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.AwsCtSqsCloudAccount.String(), retries, err,
			)
			return retry.RetryableError(fmt.Errorf(
				"Unable to create %s cloud account integration (retrying %d more time(s))",
				api.AwsCtSqsCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		d.SetId(cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)

		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type) // @afiune should we deprecate?
		d.Set("org_level", cloudAccount.IsOrg == 1)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.AwsCtSqsCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsCloudTrailRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.AwsCtSqsCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetAwsCtSqs(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	cloudAccount := response.Data
	if cloudAccount.IntgGuid == d.Id() {
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)
		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type)
		d.Set("org_level", cloudAccount.IsOrg == 1)

		creds := make(map[string]string)
		credentials := cloudAccount.Data.Credentials
		creds["role_arn"] = credentials.RoleArn
		creds["external_id"] = credentials.ExternalID
		d.Set("credentials", []map[string]string{creds})
		d.Set("queue_url", cloudAccount.Data.QueueUrl)

		accountMapFileBytes, err := cloudAccount.Data.DecodeAccountMappingFile()
		if err != nil {
			return err
		}

		accountMapFile := new(accountMappingsFile)
		if len(accountMapFileBytes) != 0 {
			// The integration has an account mapping file
			// unmarshal its content into the account mapping struct
			err := json.Unmarshal(accountMapFileBytes, accountMapFile)
			if err != nil {
				return fmt.Errorf("Error decoding organization account mapping: %s", err)
			}

		}

		err = d.Set("org_account_mappings", flattenOrgAccountMappings(accountMapFile))
		if err != nil {
			return fmt.Errorf("Error flattening organization account mapping: %s", err)
		}

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.AwsCtSqsCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework     = meta.(*api.Client)
		awsCtSqsData = api.AwsCtSqsData{
			QueueUrl: d.Get("queue_url").(string),
			Credentials: api.AwsCtSqsCredentials{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
		}
	)

	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d)
	if !accountMapFile.Empty() {
		accountMapFileBytes, err := json.Marshal(accountMapFile)
		if err != nil {
			return err
		}

		awsCtSqsData.EncodeAccountMappingFile(accountMapFileBytes)
	}

	awsCtSqs := api.NewCloudAccount(d.Get("name").(string),
		api.AwsCtSqsCloudAccount,
		awsCtSqsData,
	)

	if !d.Get("enabled").(bool) {
		awsCtSqs.Enabled = 0
	}

	awsCtSqs.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsCtSqsCloudAccount.String(), awsCtSqs.IntgGuid)
	response, err := lacework.V2.CloudAccounts.UpdateAwsCtSqs(awsCtSqs)
	if err != nil {
		return err
	}

	cloudAccount := response.Data
	d.Set("name", cloudAccount.Name)
	d.Set("intg_guid", cloudAccount.IntgGuid)
	d.Set("enabled", cloudAccount.Enabled == 1)
	d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
	d.Set("type_name", cloudAccount.Type)
	d.Set("org_level", cloudAccount.IsOrg == 1)

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.AwsCtSqsCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.AwsCtSqsCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.AwsCtSqsCloudAccount.String(), d.Id())
	return nil
}
