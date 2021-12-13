package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelEmailCreate applies integration terraform:
// => '../examples/resource_lacework_alert_channel_email'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new alert channel name and destroys it
func TestAlertChannelEmailCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_email",
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Email Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "Email Alert Channel Example", GetIntegrationName(create, "EMAIL_USER"))

	// Update Email Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Email Alert Channel Updated"}

	update := terraform.Apply(t, terraformOptions)
	assert.Equal(t, "Email Alert Channel Updated", GetIntegrationName(update, "EMAIL_USER"))
}
