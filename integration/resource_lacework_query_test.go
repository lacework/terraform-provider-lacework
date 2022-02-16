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
			"query_id": "Lql_Terraform_Query",
			"query":    queryString},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Query
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetQueryProps(create)

	actualQueryID := terraform.Output(t, terraformOptions, "query_id")
	actualQuery := terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", createProps.Data.QueryID)
	assert.Equal(t, queryString, createProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, queryString, actualQuery)

	// Update Query
	terraformOptions.Vars = map[string]interface{}{
		"query_id": "Lql_Terraform_Query",
		"query":    updatedQueryString,
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetQueryProps(update)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", updateProps.Data.QueryID)
	assert.Equal(t, updatedQueryString, updateProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, updatedQueryString, actualQuery)
}

var (
	queryString = `Lql_Terraform_Query {
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }
}`

	updatedQueryString = `Lql_Terraform_Query {
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and EVENT:userIdentity."type"::String not in ('IAMUser')
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }        
}`
)
