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
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_ecr",
		Vars: map[string]interface{}{
			"integration_name": "Amazon Elastic Container Registry Example",
			"role_arn":         "arn:aws:iam::1234567890:role/lacework_iam_example_role",
			"external_id":      "12345",
			"registry_domain":  "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com",
			"non_os_packages":  false,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Google Artifact Registry
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createData := GetEcrWithCrossAccountCreds(create)
	assert.Equal(t, "Amazon Elastic Container Registry Example", createData.Name)
	assert.Equal(t, "arn:aws:iam::1234567890:role/lacework_iam_example_role", createData.Data.Credentials.RoleArn)
	assert.Equal(t, "12345", createData.Data.Credentials.ExternalID)
	assert.Equal(t, "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com", createData.Data.RegistryDomain)
	assert.Equal(t, false, createData.Data.AwsEcrCommonData.NonOSPackageEval)

	// Update Google Artifact Registry
	terraformOptions.Vars["integration_name"] = "Amazon Elastic Container Registry Updated"
	terraformOptions.Vars["non_os_packages"] = true

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateData := GetEcrWithCrossAccountCreds(update)

	assert.Equal(t, "Amazon Elastic Container Registry Updated", updateData.Name)
	assert.Equal(t, "arn:aws:iam::1234567890:role/lacework_iam_example_role", updateData.Data.Credentials.RoleArn)
	assert.Equal(t, "12345", updateData.Data.Credentials.ExternalID)
	assert.Equal(t, "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com", updateData.Data.RegistryDomain)
	assert.Equal(t, true, updateData.Data.AwsEcrCommonData.NonOSPackageEval)
}
