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

	// Attempt to update query_id should return error
	terraformOptions.Vars = map[string]interface{}{
		"query_id": "Lql_Terraform_Query_Changed",
		"query":    updatedQueryString,
	}

	msg, err := terraform.ApplyE(t, terraformOptions)

	assert.Error(t, err)
	assert.Contains(t, msg, "unable to change ID of an existing query")
}

func TestQueryCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars: map[string]interface{}{
			"query_id": "Lql_Terraform_Query",
			"query":    queryStringK8},
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
		"query_id": "Lql_Terraform_Query",
		"query":    queryStringK8,
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

func TestQueryDeprecatedSytaxWithID(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_query",
		Vars: map[string]interface{}{
			"query_id": "Lql_Terraform_Query",
			"query":    queryDeprecatedSyntaxWithID},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Query
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetQueryProps(create)

	actualQueryID := terraform.Output(t, terraformOptions, "query_id")
	actualQuery := terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", createProps.Data.QueryID)
	assert.Equal(t, queryDeprecatedSyntaxWithID, createProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, queryDeprecatedSyntaxWithID, actualQuery)

	// Update Query
	terraformOptions.Vars = map[string]interface{}{
		"query_id": "Lql_Terraform_Query",
		"query":    updateQueryDeprecatedSyntaxWithID,
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetQueryProps(update)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", updateProps.Data.QueryID)
	assert.Equal(t, updateQueryDeprecatedSyntaxWithID, updateProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, updateQueryDeprecatedSyntaxWithID, actualQuery)

	// Run apply again
	thirdApply := terraform.ApplyAndIdempotent(t, terraformOptions)

	thirdApplyProps := GetQueryProps(thirdApply)

	actualQueryID = terraform.Output(t, terraformOptions, "query_id")
	actualQuery = terraform.Output(t, terraformOptions, "query")

	assert.Equal(t, "Lql_Terraform_Query", thirdApplyProps.Data.QueryID)
	assert.Equal(t, updateQueryDeprecatedSyntaxWithID, thirdApplyProps.Data.QueryText)

	assert.Equal(t, "Lql_Terraform_Query", actualQueryID)
	assert.Equal(t, updateQueryDeprecatedSyntaxWithID, actualQuery)
}

var (
	queryString = `{
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
	queryStringK8 = `{
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

	updatedQueryString = `{
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

	queryDeprecatedSyntaxWithID = `Lql_Terraform_Query {
    source {
        CloudTrailRawEvents
    }
    filter {
        ERROR_CODE is null
    }
    return distinct {
        EVENT
    }
}`

	updateQueryDeprecatedSyntaxWithID = `Lql_Terraform_Query{
    source { CloudTrailRawEvents }
    filter { ERROR_CODE is null }
    return distinct { EVENT }
}`
)
