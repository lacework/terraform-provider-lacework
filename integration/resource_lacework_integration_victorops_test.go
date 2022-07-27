package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestVictorOpsAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_victorops'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestVictorOpsAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_victorops",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"channel_name": "VictorOps Alert Channel Example",
			"webhook_url":  "https://alert.victorops.com/integrations/generic/20131114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new VictorOps Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "VictorOps Alert Channel Example", GetIntegrationName(create, "VICTOR_OPS"))

	// Update VictorOps Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "VictorOps Alert Channel Updated",
		"webhook_url":  "https://alert.victorops.com/integrations/generic/12331114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database"}

	update := terraform.Apply(t, terraformOptions)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualUrl := terraform.Output(t, terraformOptions, "webhook_url")
	assert.Equal(t, "VictorOps Alert Channel Updated", GetIntegrationName(update, "VICTOR_OPS"))
	assert.Equal(t, "VictorOps Alert Channel Updated", actualName)
	assert.Equal(t, "https://alert.victorops.com/integrations/generic/12331114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database", actualUrl)
}
