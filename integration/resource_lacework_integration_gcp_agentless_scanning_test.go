package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGcpAgentlessScanningCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gcp_agentless_scanning'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGcpAgentlessScanningCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	integration_name := "GCP Agentless Scanning Example Integration Test"
	update_integration_name := fmt.Sprintf("%s Updated", integration_name)
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcp_agentless_scanning",
			Vars: map[string]interface{}{
				"integration_name": integration_name,
				"client_id":        gcreds.ClientID,
				"client_email":     gcreds.ClientEmail,
				"private_key_id":   gcreds.PrivateKeyID,
				"bucket_name":      "storage bucket id",
			},
			EnvVars: map[string]string{
				"TF_VAR_private_key": gcreds.PrivateKey,
				"LW_API_TOKEN":       LwApiToken,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Google Agentless Scanning integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetGcpAgentlessScanningResponse(create)
		assert.Equal(t, integration_name, createData.Data.Name)

		// Update Gcp integration
		terraformOptions.Vars["integration_name"] = update_integration_name

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetGcpAgentlessScanningResponse(update)
		assert.Equal(t, update_integration_name, updateData.Data.Name)
	}
}
