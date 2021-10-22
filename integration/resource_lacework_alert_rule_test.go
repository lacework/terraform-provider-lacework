package integration

import (
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
			"description": "Alert Rule created by Terraform",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Alert Rule
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAlertRuleProps(create)
	assert.Equal(t, "Alert Rule created by Terraform", createProps.Data.Filter.Description)

	// Update Alert Rule
	terraformOptions.Vars["description"] = "Updated Alert Rule created by Terraform"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAlertRuleProps(update)
	assert.Equal(t, "Updated Alert Rule created by Terraform", updateProps.Data.Filter.Description)
}
