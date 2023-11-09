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
// nolint

func TestIntegrationAwsOrgAgentlessScanningLog(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_aws_org_agentless_scanning",
		Vars: map[string]interface{}{
			"name":        "AWS Org Agentless Scanning created by terraform",
			"role_arn":    "arn:aws:iam::249446771485:role/lacework-iam-example-role",
			"external_id": "12345",
			"org_account_mappings": []map[string]interface{}{
				{
					"default_lacework_account": "customerdemo",
					"mapping": []map[string]interface{}{
						{
							"lacework_account": "abc",
							"aws_accounts":     []string{"327958430571"},
						},
					},
				},
			},
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	// Create new AWS Agentless Scanning Integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetAwsAgentlessOrgScanningResponse(create)
	actualName := terraform.Output(t, terraformOptions, "name")
	assert.Equal(
		t,
		"AWS Org Agentless Scanning created by terraform",
		GetCloudOrgAccountIntegrationName(create),
	)
	assert.Equal(t, "AWS Org Agentless Scanning created by terraform", createData.Data.Name)
	assert.Equal(t, "AWS Org Agentless Scanning created by terraform", actualName)

	// Update AWS Agentless Scanning Integration
	terraformOptions.Vars = map[string]interface{}{
		"name":        "AWS Org Agentless Scanning updated by terraform",
		"role_arn":    "arn:aws:iam::249446771485:role/lacework-iam-example-role",
		"external_id": "12345",
		"org_account_mappings": []map[string]interface{}{
			{
				"default_lacework_account": "customerdemo",
				"mapping": []map[string]interface{}{
					{
						"lacework_account": "abc",
						"aws_accounts":     []string{"327958430571"},
					},
				},
			},
		},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetAwsAgentlessOrgScanningResponse(update)
	actualName = terraform.Output(t, terraformOptions, "name")
	assert.Equal(
		t,
		"AWS Org Agentless Scanning updated by terraform",
		GetCloudOrgAccountIntegrationName(update),
	)
	assert.Equal(t, "AWS Org Agentless Scanning updated by terraform", updateData.Data.Name)
	assert.Equal(t, "AWS Org Agentless Scanning updated by terraform", actualName)
}
