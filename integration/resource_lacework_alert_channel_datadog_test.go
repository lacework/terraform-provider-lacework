package integration

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
)

func TestDatadogAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_datadog",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}
