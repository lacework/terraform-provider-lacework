package integration

import (
	"fmt"
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

		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetEcrWithCrossAccountCreds(create)
		assert.Equal(t, "Amazon Elastic Container Registry Example", createData.Name)
		assert.Equal(t, awsCreds.RoleArn, createData.Data.Credentials.RoleArn)
		assert.Equal(t, awsCreds.ExternalID, createData.Data.Credentials.ExternalID)
		assert.Equal(t, awsCreds.RegistryDomain, createData.Data.RegistryDomain)
		assert.Equal(t, true, createData.Data.AwsEcrCommonData.NonOSPackageEval)

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

func TestIntegrationECRNumImagesValidation(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_ecr/iam_role",
		Vars: map[string]interface{}{
			"integration_name": "Amazon Elastic Container Registry Example",
			"num_images":       0,
		},
	})

	validationTests := []struct {
		numImages int
		valid     bool
	}{
		{numImages: 0, valid: false},
		{numImages: 3, valid: false},
		{numImages: 5, valid: true},
		{numImages: 6, valid: false},
		{numImages: 10, valid: true},
		{numImages: 12, valid: false},
		{numImages: 15, valid: true},
		{numImages: 16, valid: false},
	}

	for _, tests := range validationTests {
		t.Run(fmt.Sprintf("%d is a %t value for limit_num_imgs", tests.numImages, tests.valid), func(t *testing.T) {
			terraformOptions.Vars["num_images"] = tests.numImages
			_, err := terraform.PlanE(t, terraformOptions)
			if !tests.valid {
				if assert.Error(t, err) {
					assert.Contains(t,
						err.Error(),
						"expected limit_num_imgs to be one of [5 10 15]",
					)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}

	//Test omit num images
	terraformOptions = terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_integration_ecr/iam_role",
		Vars: map[string]interface{}{
			"integration_name": "Amazon Elastic Container Registry Example",
		},
	})
	_, err := terraform.PlanE(t, terraformOptions)
	assert.Nil(t, err)
}
