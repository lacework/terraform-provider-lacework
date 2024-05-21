package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAzureAdAl applies integration terraform:
// => '../examples/resource_lacework_integration_azure_ad_al'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationAzureAdAl(t *testing.T) {
	integration_name := "Azure Ad Al Example Integration Test With Terraform"
	updated_integration_name := fmt.Sprintf("%s Updated", integration_name)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_azure_ad_al",
		Vars: map[string]interface{}{
			"name": integration_name,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new AzureAdAl integration
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	intgRes := GetCloudAccountAzureAdAlIntegrationResponse(create)
	assert.Equal(t, integration_name, intgRes.Data.Name)

	// Update integration
	terraformOptions.Vars["name"] = updated_integration_name

	update := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	intgRes = GetCloudAccountAzureAdAlIntegrationResponse(update)
	assert.Equal(t, updated_integration_name, intgRes.Data.Name)
}
