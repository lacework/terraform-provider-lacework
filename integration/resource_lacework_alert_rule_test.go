package integration

import (
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"

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
		},
	})
	terraformOptions.TimeBetweenRetries = 2 * time.Second

	defer terraform.Destroy(t, terraformOptions)

	// Create new Alert Rule
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualName := terraform.Output(t, terraformOptions, "name")
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualChannels := terraform.Output(t, terraformOptions, "channels")
	actualSeverities := terraform.Output(t, terraformOptions, "severities")
	actualEventCategories := terraform.Output(t, terraformOptions, "event_categories")
	actualResourceGroupID := terraform.Output(t, terraformOptions, "resource_group_id")

	assert.Equal(t, "Alert Rule created by Terraform", createProps.Data.Filter.Description)
	assert.Equal(t, []string{"TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350"}, createProps.Data.Channels)
	assert.Equal(t, []string{"Critical"}, api.NewAlertRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{actualResourceGroupID}, createProps.Data.Filter.ResourceGroups)
	assert.Equal(t, []string{"Compliance"}, createProps.Data.Filter.EventCategories)

	assert.Equal(t, "Alert Rule", actualName)
	assert.Equal(t, "Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350]", actualChannels)
	assert.Equal(t, string("[Critical]"), actualSeverities)
	assert.Equal(t, "[Compliance]", actualEventCategories)

	// Update Alert Rule
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated Alert Rule created by Terraform",
		"channels": []string{"TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350",
			"TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53"},
		"severities":       []string{"High", "Medium"},
		"event_categories": []string{"Compliance", "User", "Platform"},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAlertRuleProps(update)
	actualName = terraform.Output(t, terraformOptions, "name")
	actualDescription = terraform.Output(t, terraformOptions, "description")
	actualChannels = terraform.Output(t, terraformOptions, "channels")
	actualSeverities = terraform.Output(t, terraformOptions, "severities")
	actualEventCategories = terraform.Output(t, terraformOptions, "event_categories")
	actualResourceGroupID = terraform.Output(t, terraformOptions, "resource_group_id")

	assert.Equal(t, "Updated Alert Rule created by Terraform", updateProps.Data.Filter.Description)
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350")
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53")
	assert.Equal(t, []string{"High", "Medium"}, api.NewAlertRuleSeveritiesFromIntSlice(updateProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{actualResourceGroupID}, createProps.Data.Filter.ResourceGroups)
	assert.Equal(t, []string{"Compliance", "User", "Platform"}, updateProps.Data.Filter.EventCategories)

	assert.Equal(t, "Alert Rule", actualName)
	assert.Equal(t, "Updated Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_5AB90986035F116604A26E1634340AC4FEDD1722A4D6A53 TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350]",
		actualChannels)
	assert.Equal(t, "[High Medium]", actualSeverities)
	assert.Equal(t, "[Compliance User Platform]", actualEventCategories)
}

func TestAlertRuleSeverities(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule",
		Vars: map[string]interface{}{
			"severities": []string{"Critical", "high", "mEdIuM", "LOW"},
		},
	})
	terraformOptions.TimeBetweenRetries = 2 * time.Second

	defer terraform.Destroy(t, terraformOptions)

	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualSeverities := terraform.Output(t, terraformOptions, "severities")

	assert.Equal(t,
		[]string{"Critical", "High", "Medium", "Low"},
		api.NewAlertRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice(),
	)
	assert.Equal(t, "[Critical High Medium Low]", actualSeverities)

	invalidOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule",
		Vars: map[string]interface{}{
			"severities": []string{"INVALID"},
		},
	})

	_, err := terraform.ApplyE(t, invalidOptions)
	if assert.Error(t, err) {
		assert.Contains(t,
			err.Error(),
			"severities.0: can only be 'Critical', 'High', 'Medium', 'Low', 'Info'",
		)
	}
}

func TestAlertRuleCategories(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule",
		Vars: map[string]interface{}{
			"event_categories": []string{"Compliance", "APP", "CloUD", "fIlE", "machine", "uSER", "PlatforM"},
		},
	})
	terraformOptions.TimeBetweenRetries = 2 * time.Second

	defer terraform.Destroy(t, terraformOptions)

	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualCategories := terraform.Output(t, terraformOptions, "event_categories")

	assert.Equal(t, []string{"Compliance", "App", "Cloud", "File", "Machine", "User", "Platform"}, createProps.Data.Filter.EventCategories)
	assert.Equal(t, "[Compliance App Cloud File Machine User Platform]", actualCategories)

	invalidOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule",
		Vars: map[string]interface{}{
			"event_categories": []string{"INVALID"},
		},
	})

	_, err := terraform.ApplyE(t, invalidOptions)
	if assert.Error(t, err) {
		assert.Contains(t,
			err.Error(),
			"event_categories.0: can only be 'Compliance', 'App', 'Cloud', 'File', 'Machine', 'User', 'Platform'",
		)
	}
}
