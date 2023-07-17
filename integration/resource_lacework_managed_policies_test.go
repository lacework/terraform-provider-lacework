package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestManagedPolicies applies integration terraform:
// => '../examples/resource_lacework_managed_policies'
//
// It uses the go-sdk to verify that the lacework managed policies can be updated correctly.
// nolint
func TestManagedPolicies(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_managed_policies",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"id_1":       "lacework-global-1",
			"enabled_1":  false,
			"severity_1": "High",
			"id_2":       "lacework-global-2",
			"enabled_2":  false,
			"severity_2": "Critical",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	policies := terraform.Output(t, terraformOptions, "policy")
	assert.Equal(t, "[map[enabled:false id:lacework-global-1 severity:high] map[enabled:false id:lacework-global-2 severity:critical]]", policies)

	policy1Props := GetPolicyPropsById("lacework-global-1")
	assert.False(t, policy1Props.Data.Enabled)
	assert.Equal(t, "high", policy1Props.Data.Severity)

	policy2Props := GetPolicyPropsById("lacework-global-2")
	assert.False(t, policy2Props.Data.Enabled)
	assert.Equal(t, "critical", policy2Props.Data.Severity)

	// Update managed policies
	terraformOptions.Vars = map[string]interface{}{
		"id_1":       "lacework-global-1",
		"enabled_1":  true,
		"severity_1": "Low",
		"id_2":       "lacework-global-2",
		"enabled_2":  true,
		"severity_2": "Medium",
	}

	terraform.ApplyAndIdempotent(t, terraformOptions)
	updatePolicy1Props := GetPolicyPropsById("lacework-global-1")
	assert.True(t, updatePolicy1Props.Data.Enabled)
	assert.Equal(t, "low", updatePolicy1Props.Data.Severity)

	updatePolicy2Props := GetPolicyPropsById("lacework-global-2")
	assert.True(t, updatePolicy2Props.Data.Enabled)
	assert.Equal(t, "medium", updatePolicy2Props.Data.Severity)
}

// TestManagedPoliciesWithDuplicateIDs applies integration terraform:
// => '../examples/resource_lacework_managed_policies'
//
// It uses the go-sdk to verify that the lacework managed policies can not have duplicate policy IDs.
// nolint
func TestManagedPoliciesWithDuplicateIDs(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_managed_policies",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"id_1":       "lacework-global-1",
			"enabled_1":  false,
			"severity_1": "High",
			"id_2":       "lacework-global-1",
			"enabled_2":  false,
			"severity_2": "Critical",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	_, err := terraform.InitAndApplyE(t, terraformOptions)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "Unable to update duplicate policy ID")
}

// TestManagedPoliciesWithCustomIDs applies integration terraform:
// => '../examples/resource_lacework_managed_policies'
//
// It uses the go-sdk to verify if the lacework managed policies can not have custom policy IDs.
// nolint
func TestManagedPoliciesWithCustomIDs(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_managed_policies",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"id_1":       "lacework-custom-1",
			"enabled_1":  false,
			"severity_1": "High",
			"id_2":       "lacework-custom-2",
			"enabled_2":  false,
			"severity_2": "Critical",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	_, err := terraform.InitAndApplyE(t, terraformOptions)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "Unable to update custom policy ID")
}
