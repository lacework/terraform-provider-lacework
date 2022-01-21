//go:build alert_channel

package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelCiscoWebexCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_cisco_webex'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelCiscoWebexCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_cisco_webex",
		Vars: map[string]interface{}{
			"channel_name": "Cisco Webex Alert Channel Example",
			"webhook_url":  "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Cisco webex Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Cisco Webex Alert Channel Example", GetIntegrationName(create, "CISCO_SPARK_WEBHOOK"))

	// Update Cisco Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Cisco Webex Alert Channel Example Updated",
		"webhook_url":  "https://webexapis.com/v1/webhooks/incoming/api-token",
	}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "Cisco Webex Alert Channel Example Updated", updateProps.Data.Name)
		assert.Equal(t, "https://webexapis.com/v1/webhooks/incoming/api-token", data["webhook"])
	}

	// Verify that the terraform resource has the correct information as expected
	actualChannelName := terraform.Output(t, terraformOptions, "channel_name")
	actualWebhookUrl := terraform.Output(t, terraformOptions, "webhook_url")
	assert.Equal(t, "Cisco Webex Alert Channel Example Updated", actualChannelName)
	assert.Equal(t, "https://webexapis.com/v1/webhooks/incoming/api-token", actualWebhookUrl)
}
