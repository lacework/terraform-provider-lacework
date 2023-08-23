package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestResourceGroupCreate applies integration terraform:
// => '../examples/resource_lacework_resource_group'
//
// It uses the go-sdk to verify the created resource group,
// applies an update with new description and destroys it
//
// Note: This test will only work in environments where RGv2 FF is enabled
func TestResourceGroupCreate(t *testing.T) {
	name := fmt.Sprintf("Terraform Test Resource Group V2 - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_resource_group",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"resource_group_name": name,
			"description":         "Terraform Test Resource Group V2",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Resource Group
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Terraform Test Resource Group V2", GetResourceGroupV2Description(create))

	// Update Resource Group
	terraformOptions.Vars["description"] = "Updated Terraform Test Resource Group V2"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	assert.Equal(t, "Updated Terraform Test Resource Group V2", GetResourceGroupV2Description(update))
}
