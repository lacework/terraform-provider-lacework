//go:build alert_channel
package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelMicrosoftTeams applies integration terraform:
// => '../examples/resource_lacework_alert_channel_microsoft_teams'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelMicrosoftTeams(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_microsoft_teams",
		Vars: map[string]interface{}{
			"channel_name": "Test Name",
			"webhook_url":  "https://outlook.office.com/webhook/api-token",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Microsoft Teams Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Test Name", GetIntegrationName(create, "MICROSOFT_TEAMS"))

	// Update Microsoft Teams Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Updated Test Name",
		"webhook_url":  "https://outlook.office.com/webhook/updated-api-token",
	}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "Updated Test Name", updateProps.Data.Name)
		assert.Equal(t, "https://outlook.office.com/webhook/updated-api-token", data["teamsUrl"])
	}

	// Verify that the terraform resource has the correct information as expected
	actualChannelName := terraform.Output(t, terraformOptions, "channel_name")
	actualWebhookUrl := terraform.Output(t, terraformOptions, "webhook_url")
	assert.Equal(t, "Updated Test Name", actualChannelName)
	assert.Equal(t, "https://outlook.office.com/webhook/updated-api-token", actualWebhookUrl)
}
