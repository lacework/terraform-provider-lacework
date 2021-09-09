package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupGcpCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_gcp'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupGcpCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_gcp",
		Vars: map[string]interface{}{
			"description": "Terraform Test Gcp Resource Group",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Gcp Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Terraform Test Gcp Resource Group", GetResourceGroupDescription(create))

	// Update Gcp Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Gcp Resource Group"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Updated Terraform Test Gcp Resource Group", GetResourceGroupDescription(update))
}
