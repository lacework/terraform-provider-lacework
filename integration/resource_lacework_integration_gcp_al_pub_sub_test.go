package integration

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGcpAlPubSub applies integration terraform:
// => '../examples/resource_lacework_integration_gcp_al_pub_sub'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGcpAlPubSub(t *testing.T) {
	subscription_id := os.Getenv("LW_PUB_SUB_SUBSCRIPTION")
	if subscription_id == "" {
		t.Log("this test requires that LW_PUB_SUB_SUBSCRIPTION is set")
		t.FailNow()
	}

	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcp_al_pub_sub",
			EnvVars: map[string]string{
				"LW_API_TOKEN":       LwApiToken,
				"TF_VAR_private_key": gcreds.PrivateKey,
			},
			Vars: map[string]interface{}{
				"name":             "GCP pub sub audit log integration example",
				"client_id":        gcreds.ClientID,
				"client_email":     gcreds.ClientEmail,
				"private_key_id":   gcreds.PrivateKeyID,
				"integration_type": "PROJECT",
				"project_id":       gcreds.ProjectID,
				"subscription_id":  subscription_id,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new GcpAlPubSub Integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		actualClientId := terraform.Output(t, terraformOptions, "client_id")
		actualClientEmail := terraform.Output(t, terraformOptions, "client_email")
		actualIntegrationType := terraform.Output(t, terraformOptions, "integration_type")
		actualProjectId := terraform.Output(t, terraformOptions, "project_id")
		actualSubscription := terraform.Output(t, terraformOptions, "subscription_id")
		assert.Equal(
			t,
			"GCP pub sub audit log integration example",
			GetCloudAccountIntegrationName(create),
		)
		assert.Equal(t, gcreds.ClientID, actualClientId)
		assert.Equal(t, gcreds.ClientEmail, actualClientEmail)
		assert.Equal(t, "PROJECT", actualIntegrationType)
		assert.Equal(t, gcreds.ProjectID, actualProjectId)
		assert.Equal(t, subscription_id, actualSubscription)

		// Get the newly created integration from the api
		createData := GetCloudAccountGcpPubSubAuditLogData(create)
		assert.Equal(t, "GCP pub sub audit log integration example", createData.Data.Name)

		// Update GcpAlPubSub Integration
		terraformOptions.Vars = map[string]interface{}{
			"name":             "GcpAlPubSub log integration updated",
			"client_id":        gcreds.ClientID,
			"client_email":     gcreds.ClientEmail,
			"private_key_id":   gcreds.PrivateKeyID,
			"integration_type": "PROJECT",
			"project_id":       gcreds.ProjectID,
			"subscription_id":  subscription_id,
		}

		terraformOptions.EnvVars = map[string]string{
			"LW_API_TOKEN":       LwApiToken,
			"TF_VAR_private_key": gcreds.PrivateKey,
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		actualClientId = terraform.Output(t, terraformOptions, "client_id")
		actualClientEmail = terraform.Output(t, terraformOptions, "client_email")
		actualIntegrationType = terraform.Output(t, terraformOptions, "integration_type")
		actualProjectId = terraform.Output(t, terraformOptions, "project_id")
		actualSubscription = terraform.Output(t, terraformOptions, "subscription_id")

		assert.Equal(
			t,
			"GcpAlPubSub log integration updated",
			GetCloudAccountIntegrationName(update),
		)
		assert.Equal(t, gcreds.ClientID, actualClientId)
		assert.Equal(t, gcreds.ClientEmail, actualClientEmail)
		assert.Equal(t, "PROJECT", actualIntegrationType)
		assert.Equal(t, gcreds.ProjectID, actualProjectId)
		assert.Equal(t, subscription_id, actualSubscription)

		// Get the newly updated integration from the api
		updateData := GetCloudAccountGcpPubSubAuditLogData(update)
		assert.Equal(t, "GcpAlPubSub log integration updated", updateData.Data.Name)

		// Update GcpAlPubSub with invalid configuration
		terraformOptions.Vars = map[string]interface{}{
			"name":             "GcpAlPubSub log integration updated",
			"client_id":        gcreds.ClientID,
			"client_email":     gcreds.ClientEmail,
			"private_key_id":   gcreds.PrivateKeyID,
			"integration_type": "ORGANIZATION",
			"project_id":       gcreds.ProjectID,
			"subscription_id":  subscription_id,
		}

		terraformOptions.EnvVars = map[string]string{
			"LW_API_TOKEN":       LwApiToken,
			"TF_VAR_private_key": gcreds.PrivateKey,
		}

		_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
		assert.Contains(t, err.Error(),
			"error updating cloud account integration: organization_id MUST be"+
				" set when integration_type is ORGANIZATION")
	}
}
