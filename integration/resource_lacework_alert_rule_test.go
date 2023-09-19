package integration

import (
	"fmt"
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
	name := fmt.Sprintf("Alert Rule - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/current",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name":                name,
			"description":         "Alert Rule created by Terraform",
			"channels":            []string{"TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76"},
			"severities":          []string{"Critical"},
			"alert_subcategories": []string{"Compliance"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraformOptions.TimeBetweenRetries = 2 * time.Second
	// Create new Alert Rule
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualName := terraform.Output(t, terraformOptions, "name")
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualChannels := terraform.Output(t, terraformOptions, "channels")
	actualSeverities := terraform.Output(t, terraformOptions, "severities")
	actualEventCategories := terraform.Output(t, terraformOptions, "alert_subcategories")
	actualResourceGroupID := terraform.Output(t, terraformOptions, "resource_group_id")

	assert.Equal(t, "Alert Rule created by Terraform", createProps.Data.Filter.Description)
	assert.Equal(t, []string{"TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76"}, createProps.Data.Channels)
	assert.Equal(t, []string{"Critical"}, api.NewAlertRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{actualResourceGroupID}, createProps.Data.Filter.ResourceGroups)
	assert.Equal(t, []string{"Compliance"}, createProps.Data.Filter.EventCategories)

	assert.Equal(t, name, actualName)
	assert.Equal(t, "Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76]", actualChannels)
	assert.Equal(t, string("[Critical]"), actualSeverities)
	assert.Equal(t, "[Compliance]", actualEventCategories)

	// Update Alert Rule
	terraformOptions.Vars = map[string]interface{}{
		"name":        name,
		"description": "Updated Alert Rule created by Terraform",
		"channels": []string{"TECHALLY_01BA9DCAF34B654254D6BF92E5C24023951C3F812B07527",
			"TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76"},
		"severities":          []string{"High", "Medium"},
		"alert_subcategories": []string{"Compliance", "User", "Platform"},
		"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAlertRuleProps(update)
	actualDescription = terraform.Output(t, terraformOptions, "description")
	actualChannels = terraform.Output(t, terraformOptions, "channels")
	actualSeverities = terraform.Output(t, terraformOptions, "severities")
	actualEventCategories = terraform.Output(t, terraformOptions, "alert_subcategories")
	actualResourceGroupID = terraform.Output(t, terraformOptions, "resource_group_id")

	assert.Equal(t, "Updated Alert Rule created by Terraform", updateProps.Data.Filter.Description)
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_01BA9DCAF34B654254D6BF92E5C24023951C3F812B07527")
	assert.Contains(t, updateProps.Data.Channels, "TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76")
	assert.Equal(t, []string{"High", "Medium"}, api.NewAlertRuleSeveritiesFromIntSlice(updateProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, []string{actualResourceGroupID}, updateProps.Data.Filter.ResourceGroups)
	assert.ElementsMatch(t, []string{"Compliance", "User", "Platform"}, updateProps.Data.Filter.EventCategories)
	assert.Equal(t, "Updated Alert Rule created by Terraform", actualDescription)
	assert.Equal(t, "[TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76 TECHALLY_01BA9DCAF34B654254D6BF92E5C24023951C3F812B07527]",
		actualChannels)
	assert.Equal(t, "[High Medium]", actualSeverities)
	assert.Equal(t, "[Compliance Platform User]", actualEventCategories)
}

func TestAlertRuleSeverities(t *testing.T) {
	name := fmt.Sprintf("Alert Rule - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/current",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name":                name,
			"severities":          []string{"Critical", "high", "mEdIuM", "LOW"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraformOptions.TimeBetweenRetries = 2 * time.Second
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualSeverities := terraform.Output(t, terraformOptions, "severities")

	assert.Equal(t,
		[]string{"Critical", "High", "Medium", "Low"},
		api.NewAlertRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice(),
	)
	assert.Equal(t, "[Critical High Medium Low]", actualSeverities)

	invalidOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/current",
		Vars: map[string]interface{}{
			"name":                name,
			"severities":          []string{"INVALID"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
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
	name := fmt.Sprintf("Alert Rule - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/current",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name": name,
			"alert_subcategories": []string{"Compliance", "App", "Cloud", "File", "Machine",
				"User", "Platform", "K8sActivity", "Registry", "SystemCall"},
			"alert_categories":    []string{"Policy"},
			"alert_sources":       []string{"AWS", "Agent"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraformOptions.TimeBetweenRetries = 2 * time.Second
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualCategories := terraform.Output(t, terraformOptions, "alert_subcategories")
	actualAlertCategories := terraform.Output(t, terraformOptions, "alert_categories")
	actualAlertSources := terraform.Output(t, terraformOptions, "alert_sources")

	assert.ElementsMatch(t, []string{"Compliance", "App", "Cloud", "File", "Machine",
		"User", "Platform", "K8sActivity", "Registry", "SystemCall"}, createProps.Data.Filter.EventCategories)
	assert.ElementsMatch(t, []string{"AWS", "Agent"}, createProps.Data.Filter.AlertSources)
	assert.ElementsMatch(t, []string{"Policy"}, createProps.Data.Filter.AlertCategories)

	assert.Equal(t, "[App Cloud Compliance File K8sActivity Machine Platform Registry SystemCall User]",
		actualCategories)
	assert.Equal(t, "[AWS Agent]", actualAlertSources)
	assert.Equal(t, "[Policy]", actualAlertCategories)

	invalidOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/current",
		Vars: map[string]interface{}{
			"name":                name,
			"alert_subcategories": []string{"INVALID"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})

	_, err := terraform.ApplyE(t, invalidOptions)
	if assert.Error(t, err) {
		assert.Contains(t,
			err.Error(),
			"expected alert_subcategories.0 to be one of [Compliance App Cloud File Machine User Platform K8sActivity Registry SystemCall]",
		)
	}
}

func TestAlertRuleDeprecatedEventCategories(t *testing.T) {
	name := fmt.Sprintf("Alert Rule - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/deprecated",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name": name,
			"event_categories": []string{"Compliance", "App", "Cloud", "File", "Machine",
				"User", "Platform", "K8sActivity", "Registry", "SystemCall"},
			"alert_categories":    []string{"Policy"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraformOptions.TimeBetweenRetries = 2 * time.Second
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)

	actualCategories := terraform.Output(t, terraformOptions, "event_categories")
	actualAlertCategories := terraform.Output(t, terraformOptions, "alert_categories")

	assert.ElementsMatch(t, []string{"Compliance", "App", "Cloud", "File", "Machine",
		"User", "Platform", "K8sActivity", "Registry", "SystemCall"}, createProps.Data.Filter.EventCategories)
	assert.ElementsMatch(t, []string{"Policy"}, createProps.Data.Filter.AlertCategories)

	assert.Equal(t, "[App Cloud Compliance File K8sActivity Machine Platform Registry SystemCall User]",
		actualCategories)
	assert.Equal(t, "[Policy]", actualAlertCategories)

	invalidOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_rule/deprecated",
		Vars: map[string]interface{}{
			"name":                name,
			"event_categories":    []string{"INVALID"},
			"resource_group_name": fmt.Sprintf("Used for Alert Rule Test - %s", time.Now()),
		},
	})

	_, err := terraform.ApplyE(t, invalidOptions)
	if assert.Error(t, err) {
		assert.Contains(t,
			err.Error(),
			"expected event_categories.0 to be one of [Compliance App Cloud File Machine User Platform K8sActivity Registry SystemCall]",
		)
	}
}
