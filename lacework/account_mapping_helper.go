package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
				"aws_accounts": castStringSlice(mapping["aws_accounts"].(*schema.Set).List()),
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
