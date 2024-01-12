package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestMetricsDataSource uses the Terraform plan at:
// => '../examples/lacework_metrics'
func TestMetricsDataSource(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data_source_lacework_metrics",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	id := terraform.Output(t, terraformOptions, "lacework_trace_id")
	assert.NotEmpty(t, id)
}
