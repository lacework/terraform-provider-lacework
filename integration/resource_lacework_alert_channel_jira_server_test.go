package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelJiraServerCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_jira_server'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelJiraServerCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_jira_server",
		Vars: map[string]interface{}{
			"channel_name":     "My Jira Server Example",
			"jira_url":         "fake-jira-url.com",
			"issue_type":       "Bug",
			"project_key":      "fake-project-key",
			"username":         "fake-username-techally",
			"password":         "fake-password",
			"group_issues_by":  "Events",
			"test_integration": false,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Jira Server Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "My Jira Server Example", GetIntegrationName(create))

	// Update Jira Server Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name":     "My Jira Server Example Updated",
		"jira_url":         "fake-jira-url-updated.com",
		"issue_type":       "Bug",
		"project_key":      "fake-project-key-updated",
		"username":         "fake-username-techally-updated",
		"password":         "fake-password",
		"group_issues_by":  "Resources",
		"test_integration": false,
	}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "My Jira Server Example Updated", updateProps.Data.Name)
		assert.Equal(t, "fake-jira-url-updated.com", data["jiraUrl"])
		assert.Equal(t, "Bug", data["issueType"])
		assert.Equal(t, "fake-project-key-updated", data["projectKey"])
		assert.Equal(t, "fake-username-techally-updated", data["username"])
		assert.Equal(t, "fake-password", data["password"])
		assert.Equal(t, "Resources", data["issueGrouping"])
		assert.Equal(t, "", data["customTemplateFile"])

		// Verify that the terraform resource has the correct information as expected
		actualChannelName := terraform.Output(t, terraformOptions, "name")
		actualJiraUrl := terraform.Output(t, terraformOptions, "jira_url")
		actualIssueType := terraform.Output(t, terraformOptions, "issue_type")
		actualProjectKey := terraform.Output(t, terraformOptions, "project_key")
		actualUsername := terraform.Output(t, terraformOptions, "username")
		actualPassword := terraform.Output(t, terraformOptions, "password")
		actualIssueGrouping := terraform.Output(t, terraformOptions, "group_issues_by")
		actualCustomTemplateFile := terraform.Output(t, terraformOptions, "custom_template_file")

		assert.Equal(t, "My Jira Server Example Updated", actualChannelName)
		assert.Equal(t, data["issueType"], actualIssueType)
		assert.Equal(t, data["jiraUrl"], actualJiraUrl)
		assert.Equal(t, data["projectKey"], actualProjectKey)
		assert.Equal(t, data["username"], actualUsername)
		assert.Equal(t, data["password"], actualPassword)
		assert.Equal(t, data["issueGrouping"], actualIssueGrouping)
		assert.Equal(t, data["customTemplateFile"], actualCustomTemplateFile)
	}
}
