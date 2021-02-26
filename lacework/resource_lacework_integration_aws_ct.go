package lacework

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

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

func resourceLaceworkIntegrationAwsCloudTrail() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationAwsCloudTrailCreate,
		Read:     resourceLaceworkIntegrationAwsCloudTrailRead,
		Update:   resourceLaceworkIntegrationAwsCloudTrailUpdate,
		Delete:   resourceLaceworkIntegrationAwsCloudTrailDelete,
		Schema:   awsCloudTrailIntegrationSchema,
		Importer: &schema.ResourceImporter{State: importLaceworkIntegration},
	}
}

func resourceLaceworkIntegrationAwsCloudTrailCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCloudTrailIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AwsIntegrationCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d)
	accountMapFileBytes, err := json.Marshal(accountMapFile)
	if err != nil {
		return err
	}

	if !accountMapFile.Empty() {
		aws.Data.EncodeAccountMappingFile(accountMapFileBytes)

		// switch this integration to be at the organization level
		aws.IsOrg = 1
	}

	// @afiune should we do this if there is sensitive information?
	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.AwsCloudTrailIntegration.String(), aws)
	response, err := lacework.Integrations.CreateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
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
		api.AwsCloudTrailIntegration.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.AwsCloudTrailIntegration.String(), d.Id())
	response, err := lacework.Integrations.GetAws(d.Id())
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
			creds["role_arn"] = integration.Data.Credentials.RoleArn
			creds["external_id"] = integration.Data.Credentials.ExternalID
			d.Set("credentials", []map[string]string{creds})
			d.Set("queue_url", integration.Data.QueueUrl)

			accountMapFileBytes, err := integration.Data.DecodeAccountMappingFile()
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

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.AwsCloudTrailIntegration.String(), integration.IntgGuid,
			)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCloudTrailIntegration,
			api.AwsIntegrationData{
				QueueUrl: d.Get("queue_url").(string),
				Credentials: api.AwsIntegrationCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	// verify if the user provided an account mapping
	accountMapFile := getResourceOrgAccountMappings(d)
	accountMapFileBytes, err := json.Marshal(accountMapFile)
	if err != nil {
		return err
	}

	if !accountMapFile.Empty() {
		aws.Data.EncodeAccountMappingFile(accountMapFileBytes)

		// switch this integration to be at the organization level
		aws.IsOrg = 1
	}

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.AwsCloudTrailIntegration.String(), aws)
	response, err := lacework.Integrations.UpdateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
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

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.AwsCloudTrailIntegration.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsCloudTrailDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.AwsCloudTrailIntegration.String(), d.Id())
	_, err := lacework.Integrations.DeleteAws(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.AwsCloudTrailIntegration.String(), d.Id())
	return nil
}

type accountMappingsFile struct {
	DefaultLaceworkAccount string                 `json:"defaultLaceworkAccountAws"`
	Mappings               map[string]interface{} `json:"integration_mappings"`
}

func (f *accountMappingsFile) Empty() bool {
	return f.DefaultLaceworkAccount == ""
}

func getResourceOrgAccountMappings(d *schema.ResourceData) *accountMappingsFile {
	accountMapFile := new(accountMappingsFile)
	accMapsInt := d.Get("org_account_mappings").([]interface{})
	if len(accMapsInt) != 0 && accMapsInt[0] != nil {
		accountMappings := accMapsInt[0].(map[string]interface{})

		accountMapFile = &accountMappingsFile{
			DefaultLaceworkAccount: accountMappings["default_lacework_account"].(string),
			Mappings:               map[string]interface{}{},
		}

		mappingSet := accountMappings["mapping"].(*schema.Set)
		for _, m := range mappingSet.List() {
			mapping := m.(map[string]interface{})
			accountMapFile.Mappings[mapping["lacework_account"].(string)] = map[string]interface{}{
				"aws_accounts": castStringArray(mapping["aws_accounts"].(*schema.Set).List()),
			}
		}

	}

	return accountMapFile
}

func flattenOrgAccountMappings(mappingFile *accountMappingsFile) []map[string]interface{} {
	orgAccMappings := make([]map[string]interface{}, 0, 1)

	if mappingFile.Empty() {
		return orgAccMappings
	}

	mappings := map[string]interface{}{
		"default_lacework_account": mappingFile.DefaultLaceworkAccount,
		"mapping":                  flattenMappings(mappingFile.Mappings),
	}

	orgAccMappings = append(orgAccMappings, mappings)
	return orgAccMappings
}

func flattenMappings(mappings map[string]interface{}) *schema.Set {
	var (
		orgAccountMappingsSchema = awsCloudTrailIntegrationSchema["org_account_mappings"].Elem.(*schema.Resource)
		mappingSchema            = orgAccountMappingsSchema.Schema["mapping"].Elem.(*schema.Resource)
		awsAccountsSchema        = mappingSchema.Schema["aws_accounts"].Elem.(*schema.Schema)
		res                      = schema.NewSet(schema.HashResource(mappingSchema), []interface{}{})
	)
	for laceworkAccount, m := range mappings {
		mappingValue := m.(map[string]interface{})
		res.Add(map[string]interface{}{
			"lacework_account": laceworkAccount,
			"aws_accounts": schema.NewSet(schema.HashSchema(awsAccountsSchema),
				mappingValue["aws_accounts"].([]interface{}),
			),
		})
	}

	return res
}
