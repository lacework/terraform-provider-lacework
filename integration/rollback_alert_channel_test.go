package integration

import (
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAlertChannelRollback use terraform code from:
// => '../examples/resource_lacework_alert_channel_microsoft_teams'
func TestAlertChannelRollback(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_microsoft_teams",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"channel_name":     "Test Name",
			"webhook_url":      "https://outlook.office.com/webhook/api-token",
			"test_integration": true,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Try to create new Alert Channel that will fail
	_, err := terraform.InitAndApplyE(t, terraformOptions)
	if assert.NotNil(t, err) {
		assertAlertChannelFailure(t, "init apply and fail", err)
	}

	// After failing once, try to create it again to verify that we clened up the Terraform state
	_, err = terraform.ApplyE(t, terraformOptions)
	if assert.NotNil(t, err) {
		assertAlertChannelFailure(t, "verify cleanup of TF state", err)
	}
}

func assertAlertChannelFailure(t *testing.T, msg string, err error) {
	t.Run(msg, func(t *testing.T) {
		// Assert the error returned to the user
		assert.Contains(t, err.Error(), "Error:")
		assert.Regexp(t,
			"\\[POST\\] https://(.*?).lacework.net/api/v2/AlertChannels/(.*?)/test",
			err.Error())
		// @afiune: This test in intermitent since our APIs sometimes return the full /test payload
		// and some other times they only return a 500 Internal Server Error
		// TODO: Investigate https://lacework.atlassian.net/browse/RAIN-28525
		//
		// assert.Contains(t,
		// err.Error(),
		// "The resource you are looking for might have been removed, had its name changed, or is temporarily unavailable")

		// Verify that we rollback the alert channel, that is, we removed it from the Lacework account
		re := regexp.MustCompile("lacework.net/api/v2/AlertChannels/(.*?)/test")
		match := re.FindStringSubmatch(err.Error())
		if assert.True(t, len(match) >= 2) {
			var response interface{}
			err := LwClient.V2.AlertChannels.Get(match[1], response)
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), "[404] Not found")
			}
		}
	})
}
