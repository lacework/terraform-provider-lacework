//go:build resource_group
package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupLwAccountCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_account'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupLwAccountCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test LwAccount Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_account",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":         "Terraform Test LwAccount Resource Group",
			"lw_accounts":         []string{"tech-ally"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new LwAccount Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetLwAccountResourceGroupProps(create)
	assert.Equal(t, "Terraform Test LwAccount Resource Group", createProps.Description)
	assert.Equal(t, []string{"tech-ally"}, createProps.LwAccounts)

	// Update LwAccount Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test LwAccount Resource Group"
	terraformOptions.Vars["lw_accounts"] = []string{"tech-ally", "mini-ally"}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetLwAccountResourceGroupProps(update)
	assert.Equal(t, "Updated Terraform Test LwAccount Resource Group", updateProps.Description)
	assert.Equal(t, []string{"tech-ally", "mini-ally"}, updateProps.LwAccounts)
}
