package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestDataExportRuleCreate applies integration terraform:
// => '../examples/resource_lacework_data_export_rule'
//
// It uses the go-sdk to verify the created data export rule,
// applies an update and destroys it
func TestDataExportRuleCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_data_export_rule",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name":             "Data Export Rule From Terraform",
			"type":             "Dataexport",
			"integration_ids":  []string{"TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"},
			"profile_versions": []string{"V1"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Data Export Rule
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetDataExportRuleProps(create)

	actualName := terraform.Output(t, terraformOptions, "name")
	actualEnabled := terraform.Output(t, terraformOptions, "enabled")
	actualType := terraform.Output(t, terraformOptions, "type")
	actualIDs := terraform.Output(t, terraformOptions, "integration_ids")
	actualProfileVersions := terraform.Output(t, terraformOptions, "profile_versions")

	assert.Equal(t, "Data Export Rule From Terraform", createProps.Data.Filter.Name)
	assert.Equal(t, []string{"V1"}, createProps.Data.Filter.ProfileVersions)
	assert.Equal(t, []string{"TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"}, createProps.Data.IDs)
	assert.Equal(t, "Dataexport", createProps.Data.Type)
	assert.Equal(t, 1, createProps.Data.Filter.Enabled)

	assert.Equal(t, "Data Export Rule From Terraform", actualName)
	assert.Equal(t, "Dataexport", actualType)
	assert.Equal(t, "[V1]", actualProfileVersions)
	assert.Equal(t, "[TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A]", actualIDs)
	assert.Equal(t, "true", actualEnabled)

	// Update Data Export Rule
	terraformOptions.Vars = map[string]interface{}{
		"name":             "Data Export Rule From Terraform Updated",
		"type":             "Dataexport",
		"integration_ids":  []string{"TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"},
		"profile_versions": []string{"V1"},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetDataExportRuleProps(update)

	actualName = terraform.Output(t, terraformOptions, "name")
	actualEnabled = terraform.Output(t, terraformOptions, "enabled")
	actualType = terraform.Output(t, terraformOptions, "type")
	actualIDs = terraform.Output(t, terraformOptions, "integration_ids")
	actualProfileVersions = terraform.Output(t, terraformOptions, "profile_versions")

	assert.Equal(t, "Data Export Rule From Terraform Updated", updateProps.Data.Filter.Name)
	assert.Equal(t, []string{"V1"}, updateProps.Data.Filter.ProfileVersions)
	assert.Equal(t, []string{"TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"}, updateProps.Data.IDs)
	assert.Equal(t, "Dataexport", updateProps.Data.Type)
	assert.Equal(t, 1, updateProps.Data.Filter.Enabled)

	assert.Equal(t, "Data Export Rule From Terraform Updated", actualName)
	assert.Equal(t, "Dataexport", actualType)
	assert.Equal(t, "[V1]", actualProfileVersions)
	assert.Equal(t, "[TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A]", actualIDs)
	assert.Equal(t, "true", actualEnabled)

}
