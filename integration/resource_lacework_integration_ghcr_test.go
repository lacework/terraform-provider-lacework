//go:build integration
package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGHCRCreate applies integration terraform:
// => '../examples/resource_lacework_integration_ghcr'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGHCRCreate(t *testing.T) {
	creds, err := ghcrLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GHCR_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_ghcr",
			Vars: map[string]interface{}{
				"username":               creds.Username,
				"ssl":                    true,
				"non_os_package_support": true,
			},
			EnvVars: map[string]string{
				"TF_VAR_password": creds.Token,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Github Container Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetContainerRegistryIntegration(create)
		assert.Equal(t, "Github Container Registry Example", createData.Name)
		assert.Equal(t, true, createData.Data.NonOSPackageEval)

		// Update Github Container Registry
		terraformOptions.Vars["integration_name"] = "Github Container Registry Updated"
		terraformOptions.Vars["non_os_package_support"] = true

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetContainerRegistryIntegration(update)
		assert.Equal(t, "Github Container Registry Updated", updateData.Name)
		assert.Equal(t, true, updateData.Data.NonOSPackageEval)
	}
}
