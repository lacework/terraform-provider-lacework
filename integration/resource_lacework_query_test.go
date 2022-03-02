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
func TestQueryCreateCloudtrail(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars:         map[string]interface{}{"query": queryString}})
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
		"query": updatedQueryString,
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetQueryProps(update)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", updateProps.Data.QueryID)
	assert.Equal(t, updatedQueryString, updateProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, updatedQueryString, actualQuery)

	// Attempt to update query id should return error
	terraformOptions.Vars = map[string]interface{}{
		"query": updatedQueryID,
	}

	msg, err := terraform.ApplyE(t, terraformOptions)
	if assert.Error(t, err) {
		assert.Contains(t, msg, "unable to change id of an existing query.")
		assert.Contains(t, msg, "Old ID: Lql_Terraform_Query")
		assert.Contains(t, msg, "New ID: Lql_Terraform_Query_Changed")
	}
}

func TestQueryCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars:         map[string]interface{}{"query": queryStringK8},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Query
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetQueryProps(create)

	actualQueryID := terraform.Output(t, terraformOptions, "query_id")
	actualQuery := terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", createProps.Data.QueryID)
	assert.Equal(t, queryStringK8, createProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, queryStringK8, actualQuery)

	// Update Query
	terraformOptions.Vars = map[string]interface{}{
		"query": queryStringK8,
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetQueryProps(update)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", updateProps.Data.QueryID)
	assert.Equal(t, queryStringK8, updateProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, queryStringK8, actualQuery)

	// Run apply again
	thirdApply := terraform.ApplyAndIdempotent(t, terraformOptions)

	thirdApplyProps := GetQueryProps(thirdApply)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", thirdApplyProps.Data.QueryID)
	assert.Equal(t, queryStringK8, thirdApplyProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, queryStringK8, actualQuery)
}

func TestQueryMalformed(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars:         map[string]interface{}{"query": malformedQuery},
	})

	msg, err := terraform.ApplyE(t, terraformOptions)
	if assert.Error(t, err) {
		assert.Contains(t, msg, "query id not found. (malformed)")
		assert.Contains(t, msg, "> Your query:")
		assert.Contains(t, msg, "{")
		assert.Contains(t, msg, "source { CloudTrailRawEvents }")
		assert.Contains(t, msg, "filter { ERROR_CODE is null }")
		assert.Contains(t, msg, "return distinct { EVENT }")
		assert.Contains(t, msg, "}")
		assert.Contains(t, msg, "> Compare provided query to the example at:")
		assert.Contains(t, msg, "https://docs.lacework.com/lql-overview")
	}
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
	queryStringK8 = `Lql_Terraform_Query {
      source {
          LW_ACT_K8S_AUDIT
      }
      filter {
          (EVENT_JSON:requestURI = '/api/v1/namespaces'
              or EVENT_JSON:requestURI like '/api/v1/namespaces?%')
          and EVENT_JSON:verb = 'create'
          and EVENT_JSON:responseStatus.code between 200 and 299
      }
      return distinct {
          EVENT_NAME,
          EVENT_OBJECT,
          CLUSTER_TYPE,
          CLUSTER_ID
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

	updatedQueryID = `Lql_Terraform_Query_Changed {
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

	// query doesn' have ID
	malformedQuery = `{
    source { CloudTrailRawEvents }
    filter { ERROR_CODE is null }
    return distinct { EVENT }
}`
)
