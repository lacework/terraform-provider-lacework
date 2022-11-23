package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGCRCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gcr'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGCRCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcr",
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

		// Create new Google Container Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetContainerRegisteryGcr(create)
		assert.Equal(t, "Google Container Registry Example", createData.Data.Name)

		assert.Contains(t, createData.Data.Data.LimitByRep, "my-repo")
		assert.Contains(t, createData.Data.Data.LimitByRep, "other-repo")

		assert.Contains(t, createData.Data.Data.LimitByTag, "dev*")
		assert.Contains(t, createData.Data.Data.LimitByTag, "*test")

		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value2"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"foo": "bar"})

		// Update Google Container Registry
		terraformOptions.Vars["integration_name"] = "Google Container Registry Updated"

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetContainerRegisteryGcr(update)
		assert.Equal(t, "Google Container Registry Updated", updateData.Data.Name)

		assert.Contains(t, createData.Data.Data.LimitByRep, "my-repo")
		assert.Contains(t, createData.Data.Data.LimitByRep, "other-repo")

		assert.Contains(t, createData.Data.Data.LimitByTag, "dev*")
		assert.Contains(t, createData.Data.Data.LimitByTag, "*test")

		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value2"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"foo": "bar"})
	}
}
