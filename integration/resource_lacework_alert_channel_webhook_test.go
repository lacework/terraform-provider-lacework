package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestWebhookAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_webhook'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestWebhookAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_webhook",
		Vars: map[string]interface{}{
			"channel_name": "Webhook Alert Channel Example",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Webhook Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Webhook Alert Channel Example", GetIntegrationName(create))

	// Update Webhook Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Webhook Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "Webhook Alert Channel Updated", GetIntegrationName(update))
}
