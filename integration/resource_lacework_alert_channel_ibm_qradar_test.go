package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIbmQRadarAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_qradar'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestIbmQRadarAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_qradar",
		Vars: map[string]interface{}{
			"channel_name":       "IbmQRadar Alert Channel Example",
			"host_url":           "https://qradar-lacework.com",
			"host_port":          4000,
			"communication_type": "HTTPS",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new IbmQRadar Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)
	data, ok := created.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualHost := terraform.Output(t, terraformOptions, "host_url")
	actualPort := terraform.Output(t, terraformOptions, "host_port")
	actualComm := terraform.Output(t, terraformOptions, "communication_type")

	assert.Equal(t, "IbmQRadar Alert Channel Example", created.Data.Name)
	assert.Equal(t, "https://qradar-lacework.com", data["qradarHostUrl"])
	assert.Equal(t, float64(4000), data["qradarHostPort"])
	assert.Equal(t, "HTTPS", data["qradarCommType"])

	assert.Equal(t, "IbmQRadar Alert Channel Example", actualName)
	assert.Equal(t, "https://qradar-lacework.com", actualHost)
	assert.Equal(t, "4000", actualPort)
	assert.Equal(t, "HTTPS", actualComm)

	// Update IbmQRadar Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name":       "IbmQRadar Alert Channel Updated",
		"host_url":           "https://qradar-lacework.com/updated",
		"host_port":          80,
		"communication_type": "HTTPS Self Signed Cert",
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)
	data, ok = updated.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName = terraform.Output(t, terraformOptions, "channel_name")
	actualHost = terraform.Output(t, terraformOptions, "host_url")
	actualPort = terraform.Output(t, terraformOptions, "host_port")
	actualComm = terraform.Output(t, terraformOptions, "communication_type")

	assert.Equal(t, "IbmQRadar Alert Channel Updated", updated.Data.Name)
	assert.Equal(t, "https://qradar-lacework.com/updated", data["qradarHostUrl"])
	assert.Equal(t, float64(80), data["qradarHostPort"])
	assert.Equal(t, "HTTPS Self Signed Cert", data["qradarCommType"])

	assert.Equal(t, "IbmQRadar Alert Channel Updated", actualName)
	assert.Equal(t, "https://qradar-lacework.com/updated", actualHost)
	assert.Equal(t, "80", actualPort)
	assert.Equal(t, "HTTPS Self Signed Cert", actualComm)
}
