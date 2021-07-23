package integration

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDatadogAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_datadog'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestDatadogAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_datadog",
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Datadog Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "Datadog Channel Alert Example", GetIntegrationName(create))

	// Update Datadog Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Datadog Channel Alert Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "Datadog Channel Alert Updated", GetIntegrationName(update))
}
