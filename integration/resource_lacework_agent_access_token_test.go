package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAgentAccessTokenCreate apply terraform:
// => '../examples/lacework_agent_access_token'
//
// It uses the go-sdk to verify the created token and destroys it
func TestAgentAccessTokenCreate(t *testing.T) {
	tokenName := fmt.Sprintf("Agent Token Terraform - %s", time.Now())
	// Create new Agent Access Token
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_agent_access_token",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"token_name": tokenName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createName := terraform.Output(t, terraformOptions, "token_name")
	assert.Equal(t, tokenName, createName)

	// Read Agent Access Token
	dataTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data_source_lacework_agent_access_token",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"token_name": tokenName,
		}})
	defer terraform.Destroy(t, dataTerraformOptions)

	terraform.InitAndApplyAndIdempotent(t, dataTerraformOptions)
	dataName := terraform.Output(t, terraformOptions, "token_name")
	assert.Equal(t, tokenName, dataName)
}
