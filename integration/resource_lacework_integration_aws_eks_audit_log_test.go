package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAwsEksAuditLog applies integration terraform:
// => '../examples/resource_lacework_integration_aws_eks_audit_log'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationAwsEksAuditLog(t *testing.T) {
	awsCreds, err := s3LoadCredentials("AWS_S3")
	if assert.Nil(t, err, "this test requires you to set AWS_S3 environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_aws_eks_audit_log",
			Vars: map[string]interface{}{
				"role_arn":    awsCreds.RoleArn,
				"external_id": awsCreds.ExternalID,
				"sns_arn":     "arn:aws:iam::249446771485:role/lacework-iam-example-role",
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new AwsEksAudit Integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "AWS EKS audit log integration example", GetIntegrationName(create, "AwsEksAudit"))

		// Update AwsEksAudit Integration
		terraformOptions.Vars = map[string]interface{}{
			"channel_name": "AwsEksAudit log integration updated",
			"role_arn":     awsCreds.RoleArn,
			"external_id":  awsCreds.ExternalID,
			"sns_arn":      "arn:aws:iam::249446771485:role/lacework-iam-example-role",
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "AwsEksAudit log integration updated", GetIntegrationName(update, "AwsEksAudit"))
	}
}
