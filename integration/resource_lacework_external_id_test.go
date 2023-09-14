package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestExternalIdCreate uses the Terraform plan at:
// => '../examples/resource_lacework_external_id'
func TestExternalIdCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_external_id",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	externalId := terraform.Output(t, terraformOptions, "external_id")
	assert.NotEmpty(t, externalId)
}
