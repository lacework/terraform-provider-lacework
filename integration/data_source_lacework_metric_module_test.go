package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestMetricDataSource uses the Terraform plan at:
// => '../examples/lacework_metric_module'
func TestMetricDataSource(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data_source_lacework_metric_module",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}
