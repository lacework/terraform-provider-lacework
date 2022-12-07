package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelGcpPubSubCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_gcp_pub_sub'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelGcpPubSubCreate(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_alert_channel_gcp_pub_sub",
			Vars: map[string]interface{}{
				"name":           "My GCP Pub Sub Example",
				"project_id":     "techally-hipstershop-275821",
				"topic_id":       "fake-topic-id",
				"issue_grouping": "Events",
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

		// Create new GCP Pub Sub Alert Channel
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "My GCP Pub Sub Example", GetAlertChannelName(create))

		// Update GCP Pub Sub Alert Channel
		terraformOptions.Vars = map[string]interface{}{
			"name":           "My GCP Pub Sub Example Updated",
			"project_id":     gcreds.ProjectID,
			"topic_id":       "fake-topic-id-updated",
			"issue_grouping": "Resources",
			"client_id":      gcreds.ClientID,
			"client_email":   gcreds.ClientEmail,
			"private_key_id": gcreds.PrivateKeyID,
		}
		terraformOptions.EnvVars = map[string]string{
			"TF_VAR_private_key": gcreds.PrivateKey,
			"LW_API_TOKEN":       LwApiToken,
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)

		// Verify that the lacework integration was created with the correct information
		updateProps := GetAlertChannelProps(update)
		if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
			assert.True(t, ok)
			assert.Equal(t, "My GCP Pub Sub Example Updated", updateProps.Data.Name)
			assert.Equal(t, gcreds.ProjectID, data["projectId"])
			assert.Equal(t, "fake-topic-id-updated", data["topicId"])
			assert.Equal(t, "Resources", data["issueGrouping"])
			assert.Equal(t, gcreds.ClientEmail, data["credentials"].(map[string]interface{})["clientEmail"])
			assert.Equal(t, gcreds.ClientID, data["credentials"].(map[string]interface{})["clientId"])

			// Verify that the terraform resource has the correct information as expected
			actualChannelName := terraform.Output(t, terraformOptions, "name")
			actualProjectID := terraform.Output(t, terraformOptions, "project_id")
			actualTopicID := terraform.Output(t, terraformOptions, "topic_id")
			actualIssueGrouping := terraform.Output(t, terraformOptions, "issue_grouping")
			actualClientId := terraform.Output(t, terraformOptions, "client_id")
			actualClientEmail := terraform.Output(t, terraformOptions, "client_email")

			assert.Equal(t, "My GCP Pub Sub Example Updated", actualChannelName)
			assert.Equal(t, gcreds.ProjectID, actualProjectID)
			assert.Equal(t, data["topicId"], actualTopicID)
			assert.Equal(t, data["issueGrouping"], actualIssueGrouping)
			assert.Equal(t, data["credentials"].(map[string]interface{})["clientId"], actualClientId)
			assert.Equal(t, data["credentials"].(map[string]interface{})["clientEmail"], actualClientEmail)
		}
	}
}
