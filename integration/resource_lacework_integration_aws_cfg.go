package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAwsCfg applies integration terraform:
// => '../examples/resource_lacework_integration_aws_cfg'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
// nolint
func TestIntegrationAwsCfg(t *testing.T) {
	awsCreds, err := awsLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set AWS_ECR_IAM environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_aws_cfg",
			Vars: map[string]interface{}{
				"name":        "AwsCfg created by terraform",
				"role_arn":    awsCreds.RoleArn,
				"external_id": awsCreds.ExternalID,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new AwsCfg Integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetCloudAccountAgentlessScanningResponse(create)
		actualName := terraform.Output(t, terraformOptions, "name")
		assert.Equal(
			t,
			"AwsCfg created by terraform",
			GetCloudAccountIntegrationName(create),
		)
		assert.Equal(t, "AwsCfg created by terraform", createData.Data.Name)
		assert.Equal(t, "AwsCfg created by terraform", actualName)

		// Update AwsCfg Integration
		terraformOptions.Vars = map[string]interface{}{
			"name":        "AwsCfg updated by terraform",
			"role_arn":    awsCreds.RoleArn,
			"external_id": awsCreds.ExternalID,
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetCloudAccountAgentlessScanningResponse(update)
		assert.Equal(
			t,
			"AwsCfg updated",
			GetCloudAccountIntegrationName(update),
		)
		assert.Equal(t, "AwsCfg updated by terraform", updateData.Data.Name)
		assert.Equal(t, "AwsCfg updated by terraform", actualName)
	}
}
