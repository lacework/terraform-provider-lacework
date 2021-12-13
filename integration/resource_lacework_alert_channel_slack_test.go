package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelSlackCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_slack'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelSlackCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_slack",
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Slack Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "Slack Alert Channel Example", GetIntegrationName(create, "SLACK_CHANNEL"))

	// Update Slack Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Slack Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "Slack Alert Channel Updated", GetIntegrationName(update, "SLACK_CHANNEL"))
}
