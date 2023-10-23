package integration

import (
	"encoding/json"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestNewRelicAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_newrelic'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestNewRelicAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_newrelic",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"channel_name": "NewRelic Insights Alert Channel Example",
			"insert_key":   "x-xx-xxxxxxxxxxxxxxxxxx",
			"account_id":   2338053,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new NewRelic Insights Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)
	data, ok := created.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName := terraform.Output(t, terraformOptions, "channel_name")
	actualInsertKey := terraform.Output(t, terraformOptions, "insert_key")
	actualAccountID := terraform.Output(t, terraformOptions, "account_id")

	assert.Equal(t, "NewRelic Insights Alert Channel Example", created.Data.Name)
	assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxx", data["insertKey"])
	assert.Equal(t, json.Number("2338053"), data["accountId"])

	assert.Equal(t, "NewRelic Insights Alert Channel Example", actualName)
	assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxx", actualInsertKey)
	assert.Equal(t, "2.338053e+06", actualAccountID)

	// Update NewRelic Insights Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "NewRelic Insights Alert Channel Updated",
		"insert_key":   "x-xx-xxxxxxxxxxxxxxxxxy",
		"account_id":   1338053,
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)
	data, ok = updated.Data.Data.(map[string]interface{})
	assert.True(t, ok)

	actualName = terraform.Output(t, terraformOptions, "channel_name")
	actualInsertKey = terraform.Output(t, terraformOptions, "insert_key")
	actualAccountID = terraform.Output(t, terraformOptions, "account_id")

	assert.Equal(t, "NewRelic Insights Alert Channel Updated", updated.Data.Name)
	assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxy", data["insertKey"])
	assert.Equal(t, json.Number("1338053"), data["accountId"])

	assert.Equal(t, "NewRelic Insights Alert Channel Updated", actualName)
	assert.Equal(t, "x-xx-xxxxxxxxxxxxxxxxxy", actualInsertKey)
	assert.Equal(t, "1.338053e+06", actualAccountID)
}
