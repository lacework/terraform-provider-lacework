package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGcpCfgCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gcp_cfg'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGcpCfgCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcp_cfg",
			Vars: map[string]interface{}{
				"client_id":      gcreds.ClientID,
				"client_email":   gcreds.ClientEmail,
				"private_key_id": gcreds.PrivateKeyID,
			},
			EnvVars: map[string]string{
				"TF_VAR_private_key": gcreds.PrivateKey,
				"LW_API_TOKEN":       LwApiToken,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Google Cfg
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetGcpCfgIntegration(create)
		assert.Equal(t, "Google Cfg Example", createData.Data.Name)

		// Update Google Artifact Registry
		terraformOptions.Vars["integration_name"] = "Google Cfg Updated"

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetContainerRegistryIntegration(update)
		assert.Equal(t, "Google Cfg Updated", updateData.Name)
	}
}
