package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type accountMappingsFile struct {
	DefaultLaceworkAccount string                 `json:"defaultLaceworkAccount"`
	Mappings               map[string]interface{} `json:"integration_mappings"`
}

func (f *accountMappingsFile) Empty() bool {
	return f.DefaultLaceworkAccount == ""
}

var awsMappingType string = "aws_accounts"
var gcpMappingType string = "gcp_projects"

func getResourceOrgAccountMappings(d *schema.ResourceData, mappingsType string) *accountMappingsFile {
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
			if mappingsType == "gcp_projects" {
				accountMapFile.Mappings[mapping["lacework_account"].(string)] = map[string]interface{}{
					"gcp_projects": castStringSlice(mapping[mappingsType].(*schema.Set).List()),
				}
			} else {
				accountMapFile.Mappings[mapping["lacework_account"].(string)] = map[string]interface{}{
					"aws_accounts": castStringSlice(mapping[mappingsType].(*schema.Set).List()),
				}
			}
		}

	}
	return accountMapFile
}

func flattenOrgAccountMappings(mappingFile *accountMappingsFile, mappingsType string) []map[string]interface{} {
	orgAccMappings := make([]map[string]interface{}, 0, 1)

	if mappingFile.Empty() {
		return orgAccMappings
	}

	mappings := map[string]interface{}{
		"default_lacework_account": mappingFile.DefaultLaceworkAccount,
		"mapping":                  flattenMappings(mappingFile.Mappings, mappingsType),
	}

	orgAccMappings = append(orgAccMappings, mappings)
	return orgAccMappings
}

func flattenMappings(mappings map[string]interface{}, mappingsType string) *schema.Set {
	var (
		awsOrgAccountMappingsSchema = awsCloudTrailIntegrationSchema["org_account_mappings"].Elem.(*schema.Resource)
		awsMappingSchema            = awsOrgAccountMappingsSchema.Schema["mapping"].Elem.(*schema.Resource)
		awsAccountsSchema           = awsMappingSchema.Schema[mappingsType].Elem.(*schema.Schema)
		awsRes                      = schema.NewSet(schema.HashResource(awsMappingSchema), []interface{}{})
	)

	for laceworkAccount, m := range mappings {
		mappingValue := m.(map[string]interface{})
		awsRes.Add(map[string]interface{}{
			"lacework_account": laceworkAccount,
			"aws_accounts": schema.NewSet(schema.HashSchema(awsAccountsSchema),
				mappingValue["aws_accounts"].([]interface{}),
			),
		})
	}

	return awsRes
}

func flattenOrgGcpAccountMappings(mappingFile *accountMappingsFile) []map[string]interface{} {
	orgAccMappings := make([]map[string]interface{}, 0, 1)

	if mappingFile.Empty() {
		return orgAccMappings
	}

	mappings := map[string]interface{}{
		"default_lacework_account": mappingFile.DefaultLaceworkAccount,
		"mapping":                  flattenGcpMappings(mappingFile.Mappings),
	}

	orgAccMappings = append(orgAccMappings, mappings)
	return orgAccMappings
}

func flattenGcpMappings(mappings map[string]interface{}) *schema.Set {
	var (
		gcpOrgAccountMappingsSchema = gcpAgentlessScanningIntegrationSchema["org_account_mappings"].Elem.(*schema.Resource)
		gcpMappingSchema            = gcpOrgAccountMappingsSchema.Schema["mapping"].Elem.(*schema.Resource)
		gcpAccountsSchema           = gcpMappingSchema.Schema["mapping"].Elem.(*schema.Schema)
		gcpRes                      = schema.NewSet(schema.HashResource(gcpMappingSchema), []interface{}{})
	)

	for laceworkAccount, m := range mappings {
		mappingValue := m.(map[string]interface{})
		gcpRes.Add(map[string]interface{}{
			"lacework_account": laceworkAccount,
			"gcp_projects": schema.NewSet(schema.HashSchema(gcpAccountsSchema),
				mappingValue["gcp_projects"].([]interface{}),
			),
		})
	}
	return gcpRes
}
