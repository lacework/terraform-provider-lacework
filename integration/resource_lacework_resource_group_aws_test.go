package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_aws'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_aws",
		Vars: map[string]interface{}{
			"description": "Terraform Test Aws Resource Group",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Aws Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Terraform Test Aws Resource Group", GetResourceGroupDescription(create))

	// Update Aws Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Aws Resource Group"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Updated Terraform Test Aws Resource Group", GetResourceGroupDescription(update))
}
