package integration

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatadogAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_datadog",
	})
	defer terraform.Destroy(t, terraformOptions)

	res := terraform.InitAndApply(t, terraformOptions)
	fmt.Printf("Res: %v", res)
	assert.Equal(t, res, "lacework_alert_channel_datadog.example: Creating...")
}