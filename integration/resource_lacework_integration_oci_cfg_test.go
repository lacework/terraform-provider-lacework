package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationOciCfg applies integration terraform:
// => '../examples/resource_lacework_integration_oci_cfg'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationOciCfg(t *testing.T) {
	ociCreds, err := ociLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set OCI_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_oci_cfg",
			Vars: map[string]interface{}{
				"name":        "OciCfg_created_by_terraform",
				"fingerprint": ociCreds.Fingerprint,
				"home_region": ociCreds.Region,
				"tenant_id":   ociCreds.TenacyID,
				"tenant_name": ociCreds.TenacyName,
				"user_ocid":   ociCreds.User,
			},
			EnvVars: map[string]string{
				"TF_VAR_private_key": ociCreds.PrivateKey,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new OciCfg Integration
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetCloudAccountAgentlessScanningResponse(create)
		actualName := terraform.Output(t, terraformOptions, "name")
		assert.Equal(
			t,
			"OciCfg_created_by_terraform",
			GetCloudAccountIntegrationName(create),
		)
		assert.Equal(t, "OciCfg_created_by_terraform", createData.Data.Name)
		assert.Equal(t, "OciCfg_created_by_terraform", actualName)

		// Update OciCfg Integration
		terraformOptions.Vars = map[string]interface{}{
			"name":        "OciCfg_updated_by_terraform",
			"fingerprint": ociCreds.Fingerprint,
			"private_key": ociCreds.PrivateKey,
			"home_region": ociCreds.Region,
			"tenant_id":   ociCreds.TenacyID,
			"tenant_name": ociCreds.TenacyName,
			"user_ocid":   ociCreds.User,
		}

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetCloudAccountAgentlessScanningResponse(update)
		actualName = terraform.Output(t, terraformOptions, "name")
		assert.Equal(
			t,
			"OciCfg_updated_by_terraform",
			GetCloudAccountIntegrationName(update),
		)
		assert.Equal(t, "OciCfg_updated_by_terraform", updateData.Data.Name)
		assert.Equal(t, "OciCfg_updated_by_terraform", actualName)
	}
}
