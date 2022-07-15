package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAwsAgentlessScanningLog applies integration terraform:
// => '../examples/resource_lacework_integration_aws_agentless_scanning'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func _TestIntegrationAwsAgentlessScanningLog(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_aws_agentless_scanning",
		Vars: map[string]interface{}{
			"name": "AWS Agentless Scanning created by terraform",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new AWS Agentless Scanning Integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetCloudAccountAgentlessScanningResponse(create)
	actualName := terraform.Output(t, terraformOptions, "name")
	assert.Equal(
		t,
		"AWS Agentless Scanning created by terraform",
		GetCloudAccountIntegrationName(create),
	)
	assert.Equal(t, "AWS Agentless Scanning created by terraform", createData.Data.Name)
	assert.Equal(t, "AWS Agentless Scanning created by terraform", actualName)

	// Update AWS Agentless Scanning Integration
	terraformOptions.Vars = map[string]interface{}{
		"name": "AWS Agentless Scanning updated by terraform",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetCloudAccountAgentlessScanningResponse(update)
	assert.Equal(
		t,
		"AWS Agentless Scanning updated",
		GetCloudAccountIntegrationName(update),
	)
	assert.Equal(t, "AWS Agentless Scanning updated by terraform", updateData.Data.Name)
	assert.Equal(t, "AWS Agentless Scanning updated by terraform", actualName)
}
