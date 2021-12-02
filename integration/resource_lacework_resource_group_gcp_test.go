package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupGcpCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group_gcp'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
func TestResourceGroupGcpCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test LwAccount Resource Group - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group_gcp",
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":  "Terraform Test Gcp Resource Group",
			"organization": "MyGcpOrg",
			"projects":     []string{"pro-123", "pro-321"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Gcp Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetGcpResourceGroupProps(create)
	assert.Equal(t, "Terraform Test Gcp Resource Group", createProps.Description)
	assert.Equal(t, "MyGcpOrg", createProps.Organization)
	assert.Equal(t, []string{"pro-123", "pro-321"}, createProps.Projects)

	// Update Gcp Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Gcp Resource Group"
	terraformOptions.Vars["organization"] = "MyGcpOrgUpdated"
	terraformOptions.Vars["projects"] = []string{"pro-123-updated", "pro-321-updated"}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetGcpResourceGroupProps(update)
	assert.Equal(t, "Updated Terraform Test Gcp Resource Group", updateProps.Description)
	assert.Equal(t, "MyGcpOrgUpdated", updateProps.Organization)
	assert.Equal(t, []string{"pro-123-updated", "pro-321-updated"}, updateProps.Projects)
}
