package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationInlineScannerCreate applies integration terraform:
// => '../examples/resource_lacework_integration_inline_scanner'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationInlineScannerCreate(t *testing.T) {
	integrationName := fmt.Sprintf("Inline Scanner Container Registry - %s", time.Now())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_inline_scanner",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name": integrationName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Inline Scanner Container Registry
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetContainerRegisteryInlineScanner(create)
	assert.Equal(t, integrationName, createData.Data.Name)
	assert.Equal(t, []map[string]string{{"foo": "bar"}}, createData.Data.Data.IdentifierTag)
	assert.Equal(t, "60", createData.Data.Data.LimitNumScan)

	// Update Inline Scanner Container Registry
	terraformOptions.Vars["name"] = "Inline Scanner Container Registry Updated"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetContainerRegisteryInlineScanner(update)
	assert.Equal(t, "Inline Scanner Container Registry Updated", updateData.Data.Name)
	assert.Equal(t, []map[string]string{{"foo": "bar"}}, createData.Data.Data.IdentifierTag)
	assert.Equal(t, "60", createData.Data.Data.LimitNumScan)

	server_token := terraform.Output(t, terraformOptions, "server_token")
	assert.NotEmpty(t, server_token)
}
