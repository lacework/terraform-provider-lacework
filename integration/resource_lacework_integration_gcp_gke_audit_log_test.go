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
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_gcp_gke_audit_log",
		Vars: map[string]interface{}{
			"client_id":        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"client_email":     "email@some-project-name.iam.gserviceaccount.com",
			"private_key_id":   "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"private_key":      "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"integration_type": "PROJECT",
			"project_id":       "example_project_id",
			"subscription":     "projects/example-project-id/subscriptions/example-subscription",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new GcpGkeAudit Integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	actualClientId := terraform.Output(t, terraformOptions, "client_id")
	actualClientEmail := terraform.Output(t, terraformOptions, "client_email")
	actualPrivateKey := terraform.Output(t, terraformOptions, "private_key")
	actualPrivateKeyId := terraform.Output(t, terraformOptions, "private_key_id")
	actualIntegrationType := terraform.Output(t, terraformOptions, "integration_type")
	actualProjectId := terraform.Output(t, terraformOptions, "project_id")
	actualSubscription := terraform.Output(t, terraformOptions, "subscription")
	assert.Equal(
		t,
		"GCP GKE audit log integration example",
		GetCloudAccountIntegrationName(create),
	)
	assert.Equal(
		t,
		"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		actualClientId,
	)
	assert.Equal(
		t,
		"email@some-project-name.iam.gserviceaccount.com",
		actualClientEmail,
	)
	assert.Equal(t, "", actualPrivateKey)
	assert.Equal(t, "", actualPrivateKeyId)
	assert.Equal(t, "PROJECT", actualIntegrationType)
	assert.Equal(t, "example_project_id", actualProjectId)
	assert.Equal(t, "projects/example-project-id/subscriptions/example-subscription", actualSubscription)

	// Update GcpGkeAudit Integration
	terraformOptions.Vars = map[string]interface{}{
		"name":             "GcpGkeAudit log integration updated",
		"client_id":        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"client_email":     "email@some-project-name.iam.gserviceaccount.com",
		"private_key_id":   "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"private_key":      "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"integration_type": "PROJECT",
		"project_id":       "example_project_id",
		"subscription":     "projects/example-project-id/subscriptions/example-subscription",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	actualClientId = terraform.Output(t, terraformOptions, "client_id")
	actualClientEmail = terraform.Output(t, terraformOptions, "client_email")
	actualPrivateKey = terraform.Output(t, terraformOptions, "private_key")
	actualPrivateKeyId = terraform.Output(t, terraformOptions, "private_key_id")
	actualIntegrationType = terraform.Output(t, terraformOptions, "integration_type")
	actualProjectId = terraform.Output(t, terraformOptions, "project_id")
	actualSubscription = terraform.Output(t, terraformOptions, "subscription")
	assert.Equal(
		t,
		"GcpGkeAudit log integration updated",
		GetCloudAccountIntegrationName(update),
	)
	assert.Equal(
		t,
		"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		actualClientId,
	)
	assert.Equal(
		t,
		"email@some-project-name.iam.gserviceaccount.com",
		actualClientEmail,
	)
	assert.Equal(t, "", actualPrivateKey)
	assert.Equal(t, "", actualPrivateKeyId)
	assert.Equal(t, "PROJECT", actualIntegrationType)
	assert.Equal(t, "example_project_id", actualProjectId)
	assert.Equal(t, "projects/example-project-id/subscriptions/example-subscription", actualSubscription)
}
