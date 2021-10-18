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
			"channel":      "Splunk Channel",
			"hec_token":    "BA696D5E-CA2F-4347-97CB-3C89F834816F",
			"host":         "host",
			"ssl":          true,
			"port":         80,
			"index":        "index",
			"_source":      "source",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Splunk Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)
	data, ok := created.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualToken := terraform.Output(t, terraformOptions, "hec_token")
	actualChannel := terraform.Output(t, terraformOptions, "channel")
	actualHost := terraform.Output(t, terraformOptions, "host")
	actualPort := terraform.Output(t, terraformOptions, "port")
	actualSsl := terraform.Output(t, terraformOptions, "ssl")

	assert.Equal(t, "Splunk Alert Channel Example", created.Data.Name)
	assert.Equal(t, "BA696D5E-CA2F-4347-97CB-3C89F834816F", data["hecToken"])
	assert.Equal(t, "Splunk Channel", data["channel"])
	assert.Equal(t, "host", data["host"])
	assert.Equal(t, float64(80), data["port"])
	assert.Equal(t, true, data["ssl"])

	assert.Equal(t, "Splunk Alert Channel Example", actualName)
	assert.Equal(t, "BA696D5E-CA2F-4347-97CB-3C89F834816F", actualToken)
	assert.Equal(t, "Splunk Channel", actualChannel)
	assert.Equal(t, "host", actualHost)
	assert.Equal(t, "80", actualPort)
	assert.Equal(t, "true", actualSsl)

	// Update Splunk Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Splunk Alert Channel Updated",
		"hec_token":    "BA696D5E-CA2F-4347-97CB-3C89F834815B",
		"channel":      "Updated Splunk Channel",
		"host":         "updated-host",
		"port":         8080,
		"ssl":          false,
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)
	data, ok = updated.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName = terraform.Output(t, terraformOptions, "channel_name")
	actualToken = terraform.Output(t, terraformOptions, "hec_token")
	actualChannel = terraform.Output(t, terraformOptions, "channel")
	actualHost = terraform.Output(t, terraformOptions, "host")
	actualPort = terraform.Output(t, terraformOptions, "port")
	actualSsl = terraform.Output(t, terraformOptions, "ssl")

	assert.Equal(t, "Splunk Alert Channel Updated", updated.Data.Name)
	assert.Equal(t, "BA696D5E-CA2F-4347-97CB-3C89F834815B", data["hecToken"])
	assert.Equal(t, "Updated Splunk Channel", data["channel"])
	assert.Equal(t, "updated-host", data["host"])
	assert.Equal(t, float64(8080), data["port"])
	assert.Equal(t, false, data["ssl"])

	assert.Equal(t, "Splunk Alert Channel Updated", actualName)
	assert.Equal(t, "BA696D5E-CA2F-4347-97CB-3C89F834815B", actualToken)
	assert.Equal(t, "Updated Splunk Channel", actualChannel)
	assert.Equal(t, "updated-host", actualHost)
	assert.Equal(t, "8080", actualPort)
	assert.Equal(t, "false", actualSsl)
}
