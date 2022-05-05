package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertProfileCreate applies integration terraform:
// => '../examples/resource_lacework_alert_profile'
//
// It uses the go-sdk to verify the created alert profile,
// applies an update and destroys it
func TestAlertProfileCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_profile",
		Vars: map[string]interface{}{
			"name":    "CUSTOM_PROFILE_TERRAFORM_TEST",
			"extends": "LW_CFG_GCP_DEFAULT_PROFILE",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Alert Profile
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertProfileProps(create)

	actualID := terraform.Output(t, terraformOptions, "name")
	actualExtends := terraform.Output(t, terraformOptions, "extends")
	actualDescription := terraform.Output(t, terraformOptions, "alert_description")

	assert.Equal(t, "LW_CFG_GCP_DEFAULT_PROFILE", createProps.Data.Extends)
	assert.Equal(t, "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} of organization {{ORGANIZATION}} region {{RESOURCE_REGION}}", createProps.Data.Alerts[0].Description)
	assert.Equal(t, "Violation", createProps.Data.Alerts[0].Name)
	assert.Equal(t, "LW Configuration GCP Violation Alert", createProps.Data.Alerts[0].EventName)
	assert.Equal(t, "{{_OCCURRENCE}} violation detected in project {{PROJECT_ID}}", createProps.Data.Alerts[0].Subject)

	assert.Equal(t, "CUSTOM_PROFILE_TERRAFORM_TEST", actualID)
	assert.Equal(t, "LW_CFG_GCP_DEFAULT_PROFILE", actualExtends)
	assert.Equal(t, "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} of organization {{ORGANIZATION}} region {{RESOURCE_REGION}}", actualDescription)

	// Update Alert Profile
	terraformOptions.Vars = map[string]interface{}{
		"name":              "CUSTOM_PROFILE_TERRAFORM_TEST",
		"extends":           "LW_CFG_GCP_DEFAULT_PROFILE",
		"alert_description": "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} of organization {{ORGANIZATION}} region {{RESOURCE_REGION}} Updated",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAlertProfileProps(update)
	actualExtends = terraform.Output(t, terraformOptions, "extends")
	actualDescription = terraform.Output(t, terraformOptions, "alert_description")

	assert.Equal(t, "LW_CFG_GCP_DEFAULT_PROFILE", updateProps.Data.Extends)
	assert.Equal(t, "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} of organization {{ORGANIZATION}} region {{RESOURCE_REGION}} Updated", updateProps.Data.Alerts[0].Description)
	assert.Equal(t, "Violation", updateProps.Data.Alerts[0].Name)
	assert.Equal(t, "LW Configuration GCP Violation Alert", updateProps.Data.Alerts[0].EventName)
	assert.Equal(t, "{{_OCCURRENCE}} violation detected in project {{PROJECT_ID}}", updateProps.Data.Alerts[0].Subject)

	assert.Equal(t, "LW_CFG_GCP_DEFAULT_PROFILE", actualExtends)
	assert.Equal(t, "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} of organization {{ORGANIZATION}} region {{RESOURCE_REGION}} Updated", actualDescription)
}

func TestAlertProfileValidate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_profile",
		Vars: map[string]interface{}{
			"name":    "LW_PROFILE_TERRAFORM_TEST",
			"extends": "LW_CFG_GCP_DEFAULT_PROFILE",
		},
	})

	msg, err := terraform.PlanE(t, terraformOptions)

	assert.Error(t, err)
	assert.Contains(t, msg, "expected value of name to not start with any of \"LW_\"")
}
