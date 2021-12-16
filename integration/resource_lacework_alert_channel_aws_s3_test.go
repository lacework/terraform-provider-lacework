package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelAwsS3Create applies integration terraform:
// => '../examples/resource_lacework_alert_channel_aws_s3'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelAwsS3Create(t *testing.T) {
	awsCreds, err := s3LoadCredentials("AWS_S3")
	s3BucketArn := s3LoadBucketArn()
	if assert.Nil(t, err, "this test requires you to set AWS_S3 environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_alert_channel_aws_s3",
			Vars: map[string]interface{}{
				"role_arn":    awsCreds.RoleArn,
				"external_id": awsCreds.ExternalID,
			},
			EnvVars: map[string]string{
				"TF_VAR_bucket_arn": s3BucketArn,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new AwsS3 Alert Channel
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "AwsS3 Alert Channel Example", GetIntegrationName(create, "AWS_S3"))

		// Update AwsS3 Alert Channel
		terraformOptions.Vars = map[string]interface{}{
			"channel_name": "AwsS3 Alert Channel Updated",
			"role_arn":     awsCreds.RoleArn,
			"external_id":  awsCreds.ExternalID,
		}

		terraformOptions.EnvVars = map[string]string{
			"TF_VAR_bucket_arn": s3BucketArn,
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		assert.Equal(t, "AwsS3 Alert Channel Updated", GetIntegrationName(update, "AWS_S3"))
	}
}
