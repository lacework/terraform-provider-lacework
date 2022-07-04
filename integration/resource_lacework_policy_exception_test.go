package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestPolicyExceptionCreate applies integration terraform:
// => '../examples/resource_lacework_policy_exception'
//
// It uses the go-sdk to verify the created policy exception,
// applies an update and destroys it
func TestPolicyExceptionCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_policy_exception",
		Vars: map[string]interface{}{
			"policy_id":    "lacework-global-46",
			"description":  "Policy Exception Created via Terraform",
			"field_key":    "accountIds",
			"field_values": []string{"*"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Policy Exception
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualPolicyID := terraform.Output(t, terraformOptions, "policy_id")
	createProps := GetPolicyExceptionProps(create, actualPolicyID)

	assert.Contains(t, "Policy Exception Created via Terraform", createProps.Data.Description)

	assert.Equal(t, "Policy Exception Created via Terraform", actualDescription)

	// Update Policy Exception
	terraformOptions.Vars = map[string]interface{}{
		"policy_id":    "lacework-global-46",
		"description":  "Policy Exception Created via Terraform Updated",
		"field_key":    "accountIds",
		"field_values": []string{"*"},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetPolicyExceptionProps(update, actualPolicyID)

	actualDescription = terraform.Output(t, terraformOptions, "description")

	assert.Contains(t, "Policy Exception Created via Terraform Updated", updateProps.Data.Description)

	assert.Equal(t, "Policy Exception Created via Terraform Updated", actualDescription)
}
