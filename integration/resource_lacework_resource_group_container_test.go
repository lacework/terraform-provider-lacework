package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupContainerCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_container'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupContainerCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test Machine Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_container",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description": "Terraform Test Container Resource Group",
			"ctr_tags":    []string{"test-tag"},
			"ctr_key":     "test-key",
			"ctr_value":   "test-value",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Container Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetContainerResourceGroupProps(create)
	assert.Equal(t, "Terraform Test Container Resource Group", createProps.Description)
	assert.Equal(t, []string{"test-tag"}, createProps.ContainerTags)
	assert.Equal(t, []map[string]string{{"test-key": "test-value"}}, createProps.ContainerLabels)

	// Update Container Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Container Resource Group"
	terraformOptions.Vars["ctr_tags"] = []string{"updated-tag"}
	terraformOptions.Vars["ctr_value"] = "updated-value"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetContainerResourceGroupProps(update)
	assert.Equal(t, "Updated Terraform Test Container Resource Group", updateProps.Description)
	assert.Equal(t, []string{"updated-tag"}, updateProps.ContainerTags)
	assert.Equal(t, []map[string]string{{"test-key": "updated-value"}}, updateProps.ContainerLabels)
}
