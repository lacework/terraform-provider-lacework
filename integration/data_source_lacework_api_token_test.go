package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestApiTokenDataSource uses the Terraform plan at:
// => '../examples/lacework_agent_access_token'
func TestApiTokenDataSource(t *testing.T) {
	// Read API access token
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data_source_lacework_api_token",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	apiToken := terraform.Output(t, terraformOptions, "lacework_api_token")
	assert.NotEmpty(t, apiToken)
}
