package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAzureAgentlessScanningCreate applies integration terraform:
// => '../examples/resource_lacework_integration_azure_agentless_scanning'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationAzureAgentlessScanningCreate(t *testing.T) {
	credential, err := azureLoadDefaultCredentials()
	integration_name := "Azure Agentless Scanning Example Integration Test"
	update_integration_name := fmt.Sprintf("%s Updated", integration_name)
	if assert.Nil(t, err, "this test requires you to set AZURE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_azure_agentless_scanning",
			Vars: map[string]interface{}{
				"integration_name":    integration_name,
				"client_id":           credential.ClientID,
				"blob_container_name": "blob container name",
			},
			EnvVars: map[string]string{
				"TF_VAR_CLIENT_SECRET": credential.ClientSecret,
				"LW_API_TOKEN":         LwApiToken,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Azure Agentless Scanning integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetAzureAgentlessScanningResponse(create)
		assert.Equal(t, integration_name, createData.Data.Name)

		// Update Azure integration
		terraformOptions.Vars["integration_name"] = update_integration_name

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetAzureAgentlessScanningResponse(update)
		assert.Equal(t, update_integration_name, updateData.Data.Name)
	}
}
