package integration

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestResourceGroupLwAccountCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_lw_account'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupLwAccountCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_lw_account",
		Vars: map[string]interface{}{
			"description": "Terraform Test LwAccount Resource Group",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new LwAccount Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetLwAccountResourceGroupProps(create)
	assert.Equal(t, "Terraform Test LwAccount Resource Group", createProps.Description)

	// Update LwAccount Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test LwAccount Resource Group"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetLwAccountResourceGroupProps(update)
	assert.Equal(t, "Updated Terraform Test LwAccount Resource Group", updateProps.Description)
}
