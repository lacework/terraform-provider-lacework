package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestDatadogAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_datadog'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestDatadogAlertChannelCreate(t *testing.T) {
	apiKey := "vatasha-fake-dd-api-key"
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_datadog",
		Vars: map[string]interface{}{
			"channel_name":    "Datadog Alert Channel Example",
			"datadog_site":    "com",
			"datadog_service": "Logs Detail",
			"api_key":         apiKey,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Datadog Alert Channel
	create := terraform.InitAndApply(t, terraformOptions)
	assert.Equal(t, "Datadog Alert Channel Example", GetIntegrationName(create))

	// Update Datadog Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Datadog Alert Channel Updated",
		"datadog_site":    "com",
		"datadog_service": "Logs Detail",
		"api_key":      apiKey}

	update := terraform.Apply(t, terraformOptions)

	// Verify that the lacework integration was created with the correct information
	updateProps := GetAlertChannelProps(update)
	data, ok := updateProps.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, "Datadog Alert Channel Updated", updateProps.Data.Name)
	assert.Equal(t, "com", data["datadogSite"])
	assert.Equal(t, "Logs Detail", data["datadogType"])

	// Verify that the terraform resource has the correct information as expected
	actualDatadogSite := terraform.Output(t, terraformOptions, "datadog_site")
	actualDatadogService := terraform.Output(t, terraformOptions, "datadog_service")
	actualApiKey := terraform.Output(t, terraformOptions, "api_key")
	assert.Equal(t, "Datadog Alert Channel Updated", GetIntegrationName(update))
	assert.Equal(t, "com", actualDatadogSite)
	assert.Equal(t, "Logs Detail", actualDatadogService)
	assert.Equal(t, apiKey, actualApiKey)
}
