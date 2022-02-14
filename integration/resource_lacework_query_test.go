package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestQueryCreate applies integration terraform:
// => '../examples/resource_lacework_query'
//
// It uses the go-sdk to verify the created query,
// applies an update and destroys it
//nolint
func TestQueryCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars: map[string]interface{}{
			"query_id": "lql-terraform-query",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Query
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetQueryProps(create)

	actualQueryID := terraform.Output(t, terraformOptions, "query_id")

	assert.Equal(t, "lql-terraform-query", createProps.Data.QueryID)

	assert.Equal(t, "lql-terraform-query", actualQueryID)

	// Update Query
	terraformOptions.Vars = map[string]interface{}{
		"query_id": "lql-terraform-query-updated",
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetQueryProps(update)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")

	assert.Equal(t, "lql-terraform-query-updated", updateProps.Data.QueryID)

	assert.Equal(t, "lql-terraform-query-updated", actualQueryID)
}
