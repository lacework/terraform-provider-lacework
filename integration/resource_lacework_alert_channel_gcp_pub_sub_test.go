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
func TestAlertChannelGcpPubSubCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_gcp_pub_sub",
		Vars: map[string]interface{}{
			"name":           "My GCP Pub Sub Example",
			"project_id":     "fake-project-id",
			"topic_id":       "fake-topic-id",
			"issue_grouping": "Events",
			"client_id":      "fake-client-id",
			"client_email":   "vatasha.white@lacework.net",
			"private_key":    "super-private-key",
			"private_key_id": "fake-private-key-id",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Cisco webex Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "My GCP Pub Sub Example", GetIntegrationName(create))

	// Update Cisco Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"name":           "My GCP Pub Sub Example Updated",
		"project_id":     "fake-project-id-updated",
		"topid_id":       "fake-topic-id-updated",
		"issue_grouping": "Resources",
		"credentials": map[string]string{
			"cliend_id":      "fake-client-id-updated",
			"client_email":   "vatasha.white+1234@lacework.net",
			"private_key":    "super-private-key-updated",
			"private_key_id": "fake-private-key-id-updated",
		},
	}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "My GCP Pub Sub Example Updated", updateProps.Data.Name)
		assert.Equal(t, "fake-project-id-updated", data["projectId"])
		assert.Equal(t, "fake-topic-id-updated", data["topicId"])
		assert.Equal(t, "issue_grouping", data["issue_grouping"])
	}

	// Verify that the terraform resource has the correct information as expected
	actualChannelName := terraform.Output(t, terraformOptions, "channel_name")
	actualWebhookUrl := terraform.Output(t, terraformOptions, "webhook_url")
	assert.Equal(t, "Cisco Webex Alert Channel Example Updated", actualChannelName)
	assert.Equal(t, "https://webexapis.com/v1/webhooks/incoming/api-token", actualWebhookUrl)
}
