package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupAzureCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_azure'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupAzureCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test Azure Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_azure",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":   "Terraform Test Azure Resource Group",
			"tenant":        "b21aa1ab-111a-11ab-a000-11aa1111a11a",
			"subscriptions": []string{"1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Azure Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetAzureResourceGroupProps(create)
	assert.Equal(t, "Terraform Test Azure Resource Group", createProps.Description)
	assert.Equal(t, "b21aa1ab-111a-11ab-a000-11aa1111a11a", createProps.Tenant)
	assert.Equal(t, []string{"1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0"}, createProps.Subscriptions)

	// Update Azure Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Azure Resource Group"
	terraformOptions.Vars["tenant"] = "b21aa1ab-111a-11ab-a000-11aa1111a11a"
	terraformOptions.Vars["subscriptions"] = []string{"231a0b2-abc0-1ab1-1abc-1a000ab0a0a0"}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetAzureResourceGroupProps(create)
	assert.Equal(t, "Updated Terraform Test Azure Resource Group", GetResourceGroupDescription(update))
	assert.Equal(t, "b21aa1ab-111a-11ab-a000-11aa1111a11a", updateProps.Tenant)
	assert.Equal(t, []string{"231a0b2-abc0-1ab1-1abc-1a000ab0a0a0"}, updateProps.Subscriptions)
}
