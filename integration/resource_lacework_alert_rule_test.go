package integration

import (
	"github.com/lacework/go-sdk/api"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertRuleCreate applies integration terraform:
// => '../examples/resource_lacework_alert_rule'
//
// It uses the go-sdk to verify the created alert rule,
// applies an update and destroys it
func TestAlertRuleCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule",
		Vars: map[string]interface{}{
			"description":      "Alert Rule created by Terraform",
			"channels":         []string{"TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350"},
			"severities":       []string{"Critical"},
			"event_categories": []string{"Compliance"},
			"resource_groups":  []string{"TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Alert Rule
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualName := terraform.Output(t, terraformOptions, "name")
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualChannels := terraform.Output(t, terraformOptions, "channels")
	actualSeverities := terraform.Output(t, terraformOptions, "severities")
	actualEventCategories := terraform.Output(t, terraformOptions, "event_categories")
	actualResourceGroups := terraform.Output(t, terraformOptions, "resource_groups")

	assert.Equal(t, "Alert Rule created by Terraform", createProps.Data.Filter.Description)
	assert.Equal(t, []string{"TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350"}, createProps.Data.Channels)
	assert.Equal(t, []string{"Critical"}, api.NewAlertRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{"TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B"}, createProps.Data.Filter.ResourceGroups)
	assert.Equal(t, []string{"Compliance"}, createProps.Data.Filter.EventCategories)

	assert.Equal(t, "Alert Rule", actualName)
	assert.Equal(t, "Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350]", actualChannels)
	assert.Equal(t, string("[Critical]"), actualSeverities)
	assert.Equal(t, "[Compliance]", actualEventCategories)
	assert.Equal(t, "[TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B]", actualResourceGroups)

	// Update Alert Rule
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated Alert Rule created by Terraform",
		"channels": []string{"TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350",
			"TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53"},
		"severities":       []string{"High", "Medium"},
		"event_categories": []string{"Compliance", "User", "Platform"},
		"resource_groups":  []string{"TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B"},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAlertRuleProps(update)
	actualName = terraform.Output(t, terraformOptions, "name")
	actualDescription = terraform.Output(t, terraformOptions, "description")
	actualChannels = terraform.Output(t, terraformOptions, "channels")
	actualSeverities = terraform.Output(t, terraformOptions, "severities")
	actualEventCategories = terraform.Output(t, terraformOptions, "event_categories")
	actualResourceGroups = terraform.Output(t, terraformOptions, "resource_groups")

	assert.Equal(t, "Updated Alert Rule created by Terraform", updateProps.Data.Filter.Description)
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350")
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53")
	assert.Equal(t, []string{"High", "Medium"}, api.NewAlertRuleSeveritiesFromIntSlice(updateProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{"TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B"}, updateProps.Data.Filter.ResourceGroups)
	assert.Equal(t, []string{"Compliance", "User", "Platform"}, updateProps.Data.Filter.EventCategories)

	assert.Equal(t, "Alert Rule", actualName)
	assert.Equal(t, "Updated Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53 TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350]",
		actualChannels)
	assert.Equal(t, "[High Medium]", actualSeverities)
	assert.Equal(t, "[Compliance User Platform]", actualEventCategories)
	assert.Equal(t, "[TECHALLY_528AA69075E54C783DCFAB0B76BE919287639FBAF26101B]", actualResourceGroups)
}
