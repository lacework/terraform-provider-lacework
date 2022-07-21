package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGcpGkeAuditLog applies integration terraform:
// => '../examples/resource_lacework_integration_gcp_gke_audit_log'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationGcpGkeAuditLog(t *testing.T) {
	gcreds, err := googleLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_gcp_gke_audit_log",
			EnvVars:      tokenEnvVar,
			Vars: map[string]interface{}{
				"name":             "GCP GKE audit log integration example",
				"client_id":        gcreds.ClientID,
				"client_email":     gcreds.ClientEmail,
				"private_key_id":   gcreds.PrivateKeyID,
				"private_key":      gcreds.PrivateKey,
				"integration_type": "PROJECT",
				"project_id":       gcreds.ProjectID,
				"subscription":     "projects/techally-hipstershop-275821/subscriptions/gcp-gke-audit-log-subscription",
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new GcpGkeAudit Integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		actualClientId := terraform.Output(t, terraformOptions, "client_id")
		actualClientEmail := terraform.Output(t, terraformOptions, "client_email")
		actualIntegrationType := terraform.Output(t, terraformOptions, "integration_type")
		actualProjectId := terraform.Output(t, terraformOptions, "project_id")
		actualSubscription := terraform.Output(t, terraformOptions, "subscription")
		assert.Equal(
			t,
			"GCP GKE audit log integration example",
			GetCloudAccountIntegrationName(create),
		)
		assert.Equal(t, gcreds.ClientID, actualClientId)
		assert.Equal(t, gcreds.ClientEmail, actualClientEmail)
		assert.Equal(t, "PROJECT", actualIntegrationType)
		assert.Equal(t, gcreds.ProjectID, actualProjectId)
		assert.Equal(t, "projects/techally-hipstershop-275821/subscriptions/gcp-gke-audit-log-subscription", actualSubscription)

		// Update GcpGkeAudit Integration
		terraformOptions.Vars = map[string]interface{}{
			"name":             "GcpGkeAudit log integration updated",
			"client_id":        gcreds.ClientID,
			"client_email":     gcreds.ClientEmail,
			"private_key_id":   gcreds.PrivateKeyID,
			"private_key":      gcreds.PrivateKey,
			"integration_type": "PROJECT",
			"project_id":       gcreds.ProjectID,
			"subscription":     "projects/techally-hipstershop-275821/subscriptions/gcp-gke-audit-log-subscription",
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		actualClientId = terraform.Output(t, terraformOptions, "client_id")
		actualClientEmail = terraform.Output(t, terraformOptions, "client_email")
		actualIntegrationType = terraform.Output(t, terraformOptions, "integration_type")
		actualProjectId = terraform.Output(t, terraformOptions, "project_id")
		actualSubscription = terraform.Output(t, terraformOptions, "subscription")
		assert.Equal(
			t,
			"GcpGkeAudit log integration updated",
			GetCloudAccountIntegrationName(update),
		)
		assert.Equal(t, gcreds.ClientID, actualClientId)
		assert.Equal(t, gcreds.ClientEmail, actualClientEmail)
		assert.Equal(t, "PROJECT", actualIntegrationType)
		assert.Equal(t, gcreds.ProjectID, actualProjectId)
		assert.Equal(t, "projects/techally-hipstershop-275821/subscriptions/gcp-gke-audit-log-subscription", actualSubscription)

		// Update GcpGkeAudit with invalid configuration
		terraformOptions.Vars = map[string]interface{}{
			"name":             "GcpGkeAudit log integration updated",
			"client_id":        gcreds.ClientID,
			"client_email":     gcreds.ClientEmail,
			"private_key_id":   gcreds.PrivateKeyID,
			"private_key":      gcreds.PrivateKey,
			"integration_type": "ORGANIZATION",
			"project_id":       gcreds.ProjectID,
			"subscription":     "projects/techally-hipstershop-275821/subscriptions/gcp-gke-audit-log-subscription",
		}

		_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
		assert.Contains(t, err.Error(),
			"error updating cloud account integration: organization_id MUST be"+
				" set when integration_type is ORGANIZATION")
	}
}
