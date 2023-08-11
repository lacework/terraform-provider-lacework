package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGcpAgentlessOrgScanningCreateAndUpdate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	integration_name := "GCP Org Agentless Scanning Example Integration Test"
	update_integration_name := fmt.Sprintf("%s Updated", integration_name)
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcp_org_agentless_scanning",
			Vars: map[string]interface{}{
				"integration_name": integration_name,
				"client_id":        gcreds.ClientID,
				"client_email":     gcreds.ClientEmail,
				"private_key_id":   gcreds.PrivateKeyID,
				"bucket_name":      "storage bucket id",
				"org_account_mappings": []map[string]interface{}{
					{
						"default_lacework_account": "customerdemo",
						"mapping": []map[string]interface{}{
							{
								"lacework_account": "tech-ally",
								"gcp_projects":     []string{"techally-test"},
							},
						},
					},
				},
			},
			EnvVars: map[string]string{
				"TF_VAR_private_key": gcreds.PrivateKey,
				"LW_API_TOKEN":       LwApiToken,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Google Agentless Scanning integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetGcpAgentlessOrgScanningResponse(create)
		assert.Equal(t, integration_name, createData.Data.Name)

		// Update Gcp integration
		terraformOptions.Vars = map[string]interface{}{
			"integration_name": update_integration_name,
			"client_id":        gcreds.ClientID,
			"client_email":     gcreds.ClientEmail,
			"private_key_id":   gcreds.PrivateKeyID,
			"bucket_name":      "storage bucket id",
			"org_account_mappings": []map[string]interface{}{
				{
					"default_lacework_account": "customerdemo",
					"mapping": []map[string]interface{}{
						{
							"lacework_account": "abc",
							"gcp_projects":     []string{"techally-test"},
						},
					},
				},
			},
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetGcpAgentlessOrgScanningResponse(update)
		assert.Equal(t, update_integration_name, updateData.Data.Name)
	}
}
