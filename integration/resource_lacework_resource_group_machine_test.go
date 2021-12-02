package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupMachineCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_machine'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupMachineCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test Machine Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_machine",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":         "Terraform Test Machine Resource Group",
			"machine_key":         "test-key",
			"machine_value":       "test-value",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Machine Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetMachineResourceGroupProps(create)
	assert.Equal(t, "Terraform Test Machine Resource Group", createProps.Description)
	assert.Equal(t, []map[string]string{{"test-key": "test-value"}}, createProps.MachineTags)

	// Update Machine Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Machine Resource Group"
	terraformOptions.Vars["machine_value"] = "updated-value"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetMachineResourceGroupProps(update)
	assert.Equal(t, "Updated Terraform Test Machine Resource Group", updateProps.Description)
	assert.Equal(t, []map[string]string{{"test-key": "updated-value"}}, updateProps.MachineTags)
}
