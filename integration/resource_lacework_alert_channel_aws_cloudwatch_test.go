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
	eventBusArn := cloudwatchEnvVarsDefault()
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_aws_cloudwatch",
		Vars: map[string]interface{}{
			"name": "AWS Cloudwatch Alert Channel Example",
		},
		EnvVars: map[string]string{
			"TF_VAR_event_bus_arn": eventBusArn,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new CloudwatchEb Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "AWS Cloudwatch Alert Channel Example", GetIntegrationName(create, "CLOUDWATCH_EB"))

	// Update CloudwatchEb Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"name": "AWS Cloudwatch Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "AWS Cloudwatch Alert Channel Updated", GetIntegrationName(update, "CLOUDWATCH_EB"))
}
