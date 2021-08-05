package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelEmailCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gar'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestIntegrationGARCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gar",
			Vars: map[string]interface{}{
				"client_id":      gcreds.ClientID,
				"client_email":   gcreds.ClientEmail,
				"private_key_id": gcreds.PrivateKeyID,
				"private_key":    gcreds.PrivateKey,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Email Alert Channel
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "Google Artifact Registry Example", GetIntegrationName(create))

		// Update Email Alert Channel
		terraformOptions.Vars["integration_name"] = "Google Artifact Registry Updated"

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "Google Artifact Registry Updated", GetIntegrationName(update))
	}
}
