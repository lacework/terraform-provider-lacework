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
//nolint
func _TestAlertChannelAwsS3Create(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_aws_s3",
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new AwsS3 Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "AwsS3 Alert Channel Example", GetIntegrationName(create))

	// Update AwsS3 Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "AwsS3 Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "AwsS3 Alert Channel Updated", GetIntegrationName(update))
}
