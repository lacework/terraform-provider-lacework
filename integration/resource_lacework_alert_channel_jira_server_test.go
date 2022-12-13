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
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"channel_name":    "My Jira Server Example",
			"jira_url":        "test-lacework.atlassian.net",
			"issue_type":      "Bug",
			"project_key":     "fake-project-key",
			"username":        "fake-username-techally",
			"password":        "fake-password",
			"configuration":   "Unidirectional",
			"group_issues_by": "Events",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Jira Server Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "My Jira Server Example", GetAlertChannelName(create))

	// Update Jira Server Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name":    "My Jira Server Example Updated",
		"jira_url":        "updatedtest-lacework.atlassian.net",
		"issue_type":      "Task",
		"project_key":     "fake-project-key-updated",
		"username":        "fake-username-techally-updated",
		"password":        "fake-password-updated",
		"configuration":   "Bidirectional",
		"group_issues_by": "Resources",
	}

	customTemplateFile := "{\n    \"fields\": {\n        \"labels\": [\n            \"myLabel\"\n        ],\n        \"priority\":\n        {\n            \"id\": \"1\"\n        }\n    }\n}\n"
	customTemplateFileEncoded := "data:application/json;name=i.json;base64,ewogICAgImZpZWxkcyI6IHsKICAgICAgICAibGFiZWxzIjogWwogICAgICAgICAgICAibXlMYWJlbCIKICAgICAgICBdLAogICAgICAgICJwcmlvcml0eSI6CiAgICAgICAgewogICAgICAgICAgICAiaWQiOiAiMSIKICAgICAgICB9CiAgICB9Cn0K"

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	if data, ok := updateProps.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)
		assert.Equal(t, "My Jira Server Example Updated", updateProps.Data.Name)
		assert.Equal(t, "updatedtest-lacework.atlassian.net", data["jiraUrl"])
		assert.Equal(t, "Task", data["issueType"])
		assert.Equal(t, "fake-project-key-updated", data["projectId"])
		assert.Equal(t, "fake-username-techally-updated", data["username"])
		assert.Equal(t, "Bidirectional", data["bidirectionalConfig"])
		assert.Equal(t, "Resources", data["issueGrouping"])
		assert.Equal(t, customTemplateFileEncoded, data["customTemplateFile"])

		// Verify that the terraform resource has the correct information as expected
		actualChannelName := terraform.Output(t, terraformOptions, "channel_name")
		actualJiraUrl := terraform.Output(t, terraformOptions, "jira_url")
		actualIssueType := terraform.Output(t, terraformOptions, "issue_type")
		actualProjectKey := terraform.Output(t, terraformOptions, "project_key")
		actualUsername := terraform.Output(t, terraformOptions, "username")
		actualPassword := terraform.Output(t, terraformOptions, "password")
		actualConfiguration := terraform.Output(t, terraformOptions, "configuration")
		actualIssueGrouping := terraform.Output(t, terraformOptions, "group_issues_by")
		actualCustomTemplateFile := terraform.Output(t, terraformOptions, "custom_template_file")

		assert.Equal(t, "My Jira Server Example Updated", actualChannelName)
		assert.Equal(t, data["issueType"], actualIssueType)
		assert.Equal(t, data["jiraUrl"], actualJiraUrl)
		assert.Equal(t, data["projectId"], actualProjectKey)
		assert.Equal(t, data["username"], actualUsername)
		assert.Equal(t, "fake-password-updated", actualPassword)
		assert.Equal(t, data["bidirectionalConfig"], actualConfiguration)
		assert.Equal(t, data["issueGrouping"], actualIssueGrouping)
		assert.Equal(t, customTemplateFile, actualCustomTemplateFile)
	}
}
