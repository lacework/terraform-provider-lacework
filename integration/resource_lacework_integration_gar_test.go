package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGARCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gar'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGARCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gar",
			Vars: map[string]interface{}{
				"client_id":              gcreds.ClientID,
				"client_email":           gcreds.ClientEmail,
				"private_key_id":         gcreds.PrivateKeyID,
				"non_os_package_support": true,
			},
			EnvVars: map[string]string{
				"TF_VAR_private_key": gcreds.PrivateKey, // @afiune this will avoid printing secrets in our pipeline
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Google Artifact Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetContainerRegistryIntegration(create)
		assert.Equal(t, "Google Artifact Registry Example", createData.Name)
		assert.Equal(t, false, createData.Data.NonOSPackageEval)

		// Update Google Artifact Registry
		terraformOptions.Vars["integration_name"] = "Google Artifact Registry Updated"
		terraformOptions.Vars["non_os_package_support"] = true

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetContainerRegistryIntegration(update)
		assert.Equal(t, "Google Artifact Registry Updated", updateData.Name)
		assert.Equal(t, true, updateData.Data.NonOSPackageEval)
	}
}
