package integration

import (
"testing"

"github.com/gruntwork-io/terratest/modules/terraform"
"github.com/stretchr/testify/assert"
)

// TestAlertChannelCloudWatchCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_aws_cloudwatch'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelCloudWatchCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_aws_cloudwatch",
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new CloudwatchEb Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "AWS Cloudwatch Alert Channel Example", GetIntegrationName(create))

	// Update CloudwatchEb Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "AWS Cloudwatch Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "AWS Cloudwatch Alert Channel Updated", GetIntegrationName(update))
}
