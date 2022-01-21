//go:build integration
package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationECRCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gar'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationECRCreate(t *testing.T) {
	awsCreds, err := ecrLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set AWS_ECR_IAM environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_ecr/iam_role",
			Vars: map[string]interface{}{
				"integration_name":       "Amazon Elastic Container Registry Example",
				"role_arn":               awsCreds.RoleArn,
				"external_id":            awsCreds.ExternalID,
				"registry_domain":        awsCreds.RegistryDomain,
				"non_os_package_support": true,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Google Artifact Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetEcrWithCrossAccountCreds(create)
		assert.Equal(t, "Amazon Elastic Container Registry Example", createData.Name)
		assert.Equal(t, awsCreds.RoleArn, createData.Data.Credentials.RoleArn)
		assert.Equal(t, awsCreds.ExternalID, createData.Data.Credentials.ExternalID)
		assert.Equal(t, awsCreds.RegistryDomain, createData.Data.RegistryDomain)
		assert.Equal(t, true, createData.Data.AwsEcrCommonData.NonOSPackageEval)

		// Update Google Artifact Registry
		terraformOptions.Vars["integration_name"] = "Amazon Elastic Container Registry Updated"
		terraformOptions.Vars["non_os_package_support"] = true

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetEcrWithCrossAccountCreds(update)

		assert.Equal(t, "Amazon Elastic Container Registry Updated", updateData.Name)
		assert.Equal(t, awsCreds.RoleArn, updateData.Data.Credentials.RoleArn)
		assert.Equal(t, awsCreds.ExternalID, updateData.Data.Credentials.ExternalID)
		assert.Equal(t, awsCreds.RegistryDomain, updateData.Data.RegistryDomain)
		assert.Equal(t, true, updateData.Data.AwsEcrCommonData.NonOSPackageEval)
	}
}
