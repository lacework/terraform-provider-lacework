package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationProxyScannerCreate applies integration terraform:
// => '../examples/resource_lacework_integration_inline_scanner'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationProxyScannerCreate(t *testing.T) {
	integrationName := fmt.Sprintf("Proxy Scanner Container Registry - %s", time.Now())

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_proxy_scanner",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"name": integrationName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Proxy Scanner Container Registry
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetContainerRegisteryProxyScanner(create)
	assert.Equal(t, integrationName, createData.Data.Name)
	assert.Equal(t, 10, createData.Data.Data.LimitNumImg)
	assert.Equal(t, []map[string]string{{"foo": "bar"}}, createData.Data.Data.LimitByLabel)
	assert.Equal(t, []string{"dev*", "*test"}, createData.Data.Data.LimitByTag)
	assert.Equal(t, []string{"repo/my-image", "repo/other-image"}, createData.Data.Data.LimitByRep)

	// Update Proxy Scanner Container Registry
	terraformOptions.Vars["name"] = "Proxy Scanner Container Registry Updated"

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetContainerRegisteryProxyScanner(update)
	assert.Equal(t, "Proxy Scanner Container Registry Updated", updateData.Data.Name)
	assert.Equal(t, 10, createData.Data.Data.LimitNumImg)
	assert.Equal(t, []map[string]string{{"foo": "bar"}}, createData.Data.Data.LimitByLabel)
	assert.Equal(t, []string{"dev*", "*test"}, createData.Data.Data.LimitByTag)
	assert.Equal(t, []string{"repo/my-image", "repo/other-image"}, createData.Data.Data.LimitByRep)
	assert.NotEmpty(t, createData.Data.ServerToken.ServerToken)
	assert.NotEmpty(t, createData.Data.ServerToken.Uri)

	server_token := terraform.Output(t, terraformOptions, "server_token")
	assert.NotEmpty(t, server_token)
}
