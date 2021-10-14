package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestSplunkAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_splunk'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestSplunkAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_splunk",
		Vars: map[string]interface{}{
			"channel_name": "Splunk Alert Channel Example",
			"hec_token":    "BA696D5E-CA2F-4347-97CB-3C89F834816F",
			"host":         "host",
			"port":         80,
			"index":        "index",
			"_source":      "source",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Splunk Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)

	assert.Equal(t, "Splunk Alert Channel Example", created.Data.Name)

	// Update Splunk Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Splunk Alert Channel Updated",
		"hec_token":    "BA696D5E-CA2F-4347-97CB-3C89F834816A",
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualUrl := terraform.Output(t, terraformOptions, "hec_token")
	assert.Equal(t, "Splunk Alert Channel Updated", updated.Data.Name)
	assert.Equal(t, "Splunk Alert Channel Updated", actualName)
	assert.Equal(t, "BA696D5E-CA2F-4347-97CB-3C89F834816A", actualUrl)
}
