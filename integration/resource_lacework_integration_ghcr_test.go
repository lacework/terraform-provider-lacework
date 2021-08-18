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
				"username": creds.Username,
				"ssl":      true,
			},
			EnvVars: map[string]string{
				"TF_VAR_password": creds.Token,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Github Container Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "Github Container Registry Example", GetIntegrationName(create))

		// Update Github Container Registry
		terraformOptions.Vars["integration_name"] = "Github Container Registry Updated"

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "Github Container Registry Example", GetIntegrationName(update))
	}
}
