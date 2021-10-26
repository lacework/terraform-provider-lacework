package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelJiraCloudCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_jira_cloud'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelJiraCloudCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_jira_cloud",
		Vars: map[string]interface{}{
			"channel_name":     "My Jira Cloud Example",
			"jira_url":         "fake-jira-url.com",
			"issue_type":       "Bug",
			"project_key":      "fake-project-key",
			"username":         "fake-username-techally",
			"api_token":        "fake-api-token",
			"group_issues_by":  "Events",
			"test_integration": false,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Jira Cloud Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "My Jira Cloud Example", GetIntegrationName(create))

	// Update Jira Cloud Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name":     "My Jira Cloud Example Updated",
		"jira_url":         "fake-jira-url-updated.com",
		"issue_type":       "Bug",
		"project_key":      "fake-project-key-updated",
		"username":         "fake-username-techally-updated",
		"api_token":        "fake-api-token",
		"group_issues_by":  "Resources",
		"test_integration": false,
	}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "My Jira Cloud Example Updated", updateProps.Data.Name)
		assert.Equal(t, "fake-jira-url-updated.com", data["jiraUrl"])
		assert.Equal(t, "Bug", data["issueType"])
		assert.Equal(t, "fake-project-key-updated", data["projectKey"])
		assert.Equal(t, "fake-username-techally-updated", data["username"])
		assert.Equal(t, "fake-api-token", data["apiToken"])
		assert.Equal(t, "Resources", data["issueGrouping"])
		assert.Equal(t, "", data["customTemplateFile"])

		// Verify that the terraform resource has the correct information as expected
		actualChannelName := terraform.Output(t, terraformOptions, "name")
		actualJiraUrl := terraform.Output(t, terraformOptions, "jira_url")
		actualIssueType := terraform.Output(t, terraformOptions, "issue_type")
		actualProjectKey := terraform.Output(t, terraformOptions, "project_key")
		actualUsername := terraform.Output(t, terraformOptions, "username")
		actualApiToken := terraform.Output(t, terraformOptions, "api_token")
		actualIssueGrouping := terraform.Output(t, terraformOptions, "group_issues_by")
		actualCustomTemplateFile := terraform.Output(t, terraformOptions, "custom_template_file")

		assert.Equal(t, "My Jira Cloud Example Updated", actualChannelName)
		assert.Equal(t, data["issueType"], actualIssueType)
		assert.Equal(t, data["jiraUrl"], actualJiraUrl)
		assert.Equal(t, data["projectKey"], actualProjectKey)
		assert.Equal(t, data["username"], actualUsername)
		assert.Equal(t, data["apiToken"], actualApiToken)
		assert.Equal(t, data["issueGrouping"], actualIssueGrouping)
		assert.Equal(t, data["customTemplateFile"], actualCustomTemplateFile)
	}
}
