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
func TestIntegrationAwsEksAuditLogWithS3(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_aws_eks_audit_log/audit_log_with_s3_arn",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"role_arn":      "arn:aws:iam::249446771485:role/lacework-iam-example-role",
			"external_id":   "12345",
			"sns_arn":       "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks",
			"s3_bucket_arn": "arn:aws:s3:::example-bucket-name",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new AwsEksAudit Integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetCloudAccountEksAuditLogData(create)
	actualRoleArn := terraform.Output(t, terraformOptions, "role_arn")
	actualExternalId := terraform.Output(t, terraformOptions, "external_id")
	actualSnsArn := terraform.Output(t, terraformOptions, "sns_arn")
	actualS3Arn := terraform.Output(t, terraformOptions, "s3_bucket_arn")
	assert.Equal(
		t,
		"AWS EKS audit log integration example with S3",
		GetCloudAccountIntegrationName(create),
	)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", createData.Credentials.RoleArn)
	assert.Equal(t, "12345", createData.Credentials.ExternalID)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", createData.SnsArn)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", actualRoleArn)
	assert.Equal(t, "12345", actualExternalId)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", actualSnsArn)
	assert.Equal(t, "arn:aws:s3:::example-bucket-name", actualS3Arn)

	// Update AwsEksAudit Integration
	terraformOptions.Vars = map[string]interface{}{
		"name":          "AwsEksAudit log integration updated",
		"role_arn":      "arn:aws:iam::249446771485:role/lacework-iam-example-role",
		"external_id":   "12345",
		"sns_arn":       "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks",
		"s3_bucket_arn": "arn:aws:s3:::example-bucket-name",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetCloudAccountEksAuditLogData(update)
	actualRoleArn = terraform.Output(t, terraformOptions, "role_arn")
	actualExternalId = terraform.Output(t, terraformOptions, "external_id")
	actualSnsArn = terraform.Output(t, terraformOptions, "sns_arn")
	actualS3Arn = terraform.Output(t, terraformOptions, "s3_bucket_arn")
	assert.Equal(
		t,
		"AwsEksAudit log integration updated",
		GetCloudAccountIntegrationName(update),
	)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", updateData.Credentials.RoleArn)
	assert.Equal(t, "12345", updateData.Credentials.ExternalID)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", updateData.SnsArn)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", actualRoleArn)
	assert.Equal(t, "12345", actualExternalId)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", actualSnsArn)
	assert.Equal(t, "arn:aws:s3:::example-bucket-name", actualS3Arn)
}

func TestIntegrationAwsEksAuditLogWithOutS3(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_aws_eks_audit_log/audit_log_without_s3_arn",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"role_arn":    "arn:aws:iam::249446771485:role/lacework-iam-example-role",
			"external_id": "12345",
			"sns_arn":     "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new AwsEksAudit Integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetCloudAccountEksAuditLogData(create)
	actualRoleArn := terraform.Output(t, terraformOptions, "role_arn")
	actualExternalId := terraform.Output(t, terraformOptions, "external_id")
	actualSnsArn := terraform.Output(t, terraformOptions, "sns_arn")
	assert.Equal(
		t,
		"AWS EKS audit log integration example no S3",
		GetCloudAccountIntegrationName(create),
	)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", createData.Credentials.RoleArn)
	assert.Equal(t, "12345", createData.Credentials.ExternalID)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", createData.SnsArn)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", actualRoleArn)
	assert.Equal(t, "12345", actualExternalId)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", actualSnsArn)

	// Update AwsEksAudit Integration
	terraformOptions.Vars = map[string]interface{}{
		"name":        "AwsEksAudit log integration updated",
		"role_arn":    "arn:aws:iam::249446771485:role/lacework-iam-example-role",
		"external_id": "12345",
		"sns_arn":     "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetCloudAccountEksAuditLogData(update)
	actualRoleArn = terraform.Output(t, terraformOptions, "role_arn")
	actualExternalId = terraform.Output(t, terraformOptions, "external_id")
	actualSnsArn = terraform.Output(t, terraformOptions, "sns_arn")
	assert.Equal(
		t,
		"AwsEksAudit log integration updated",
		GetCloudAccountIntegrationName(update),
	)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", updateData.Credentials.RoleArn)
	assert.Equal(t, "12345", updateData.Credentials.ExternalID)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", updateData.SnsArn)
	assert.Equal(t, "arn:aws:iam::249446771485:role/lacework-iam-example-role", actualRoleArn)
	assert.Equal(t, "12345", actualExternalId)
	assert.Equal(t, "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks", actualSnsArn)
}
