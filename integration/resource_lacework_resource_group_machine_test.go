package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupMachineCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_machine'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupMachineCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_machine",
		Vars: map[string]interface{}{
			"description": "Terraform Test Machine Resource Group",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Machine Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Terraform Test Machine Resource Group", GetResourceGroupDescription(create))

	// Update Machine Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Machine Resource Group"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Updated Terraform Test Machine Resource Group", GetResourceGroupDescription(update))
}
