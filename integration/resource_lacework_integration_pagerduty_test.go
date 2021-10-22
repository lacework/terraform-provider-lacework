package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestPagerDutyAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_qradar'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestPagerDutyAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_qradar",
		Vars: map[string]interface{}{
			"channel_name":    "PagerDuty Alert Channel Example",
			"integration_key": "1234abc8901abc567abc123abc78e012",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new PagerDuty Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)
	data, ok := created.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualIntKey := terraform.Output(t, terraformOptions, "integration_key")

	assert.Equal(t, "PagerDuty Alert Channel Example", created.Data.Name)
	assert.Equal(t, "1234abc8901abc567abc123abc78e012", data["apiIntgKey"])

	assert.Equal(t, "PagerDuty Alert Channel Example", actualName)
	assert.Equal(t, "1234abc8901abc567abc123abc78e012", actualIntKey)

	// Update PagerDuty Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name":    "PagerDuty Alert Channel Updated",
		"integration_key": "1234abc8901abc567abc123abc78e013",
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)
	data, ok = updated.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName = terraform.Output(t, terraformOptions, "channel_name")
	actualIntKey = terraform.Output(t, terraformOptions, "integration_key")

	assert.Equal(t, "PagerDuty Alert Channel Updated", updated.Data.Name)
	assert.Equal(t, "1234abc8901abc567abc123abc78e013", data["apiIntgKey"])

	assert.Equal(t, "PagerDuty Alert Channel Updated", actualName)
	assert.Equal(t, "1234abc8901abc567abc123abc78e013", actualIntKey)
}
