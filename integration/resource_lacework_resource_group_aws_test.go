//go:build resource_group

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupAwsCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_aws'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupAwsCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test Aws Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_aws",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":         "Terraform Test Aws Resource Group",
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
