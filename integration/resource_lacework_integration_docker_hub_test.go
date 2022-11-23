package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationDockerhub applies integration terraform:
// => '../examples/resource_lacework_integration_docker_hub'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
func TestIntegrationDockerhub(t *testing.T) {
	creds, err := dockerLoadDefaultCredentials()
	if assert.Nil(t, err, "this test requires you to set DOCKER_CREDENTIALS environment variable") {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../examples/resource_lacework_integration_docker_hub",
			Vars: map[string]interface{}{
				"user":                   creds.Username,
				"non_os_package_support": false,
			},
			EnvVars: map[string]string{
				"TF_VAR_pass":  creds.Password,
				"LW_API_TOKEN": LwApiToken,
			},
		})
		defer terraform.Destroy(t, terraformOptions)

		// Create new Dockerhub Container Registry
		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
		createData := GetContainerRegisteryDockerhub(create)
		assert.Equal(t, "Dockerhub Container Registry Example", createData.Data.Name)
		assert.Equal(t, false, createData.Data.Data.NonOSPackageEval)

		assert.Contains(t, createData.Data.Data.LimitByRep, "my-repo")
		assert.Contains(t, createData.Data.Data.LimitByRep, "other-repo")

		assert.Contains(t, createData.Data.Data.LimitByTag, "dev*")
		assert.Contains(t, createData.Data.Data.LimitByTag, "*test")

		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"key": "value2"})
		assert.Contains(t, createData.Data.Data.LimitByLabel, map[string]string{"foo": "bar"})

		// Update Dockerhub Container Registry
		terraformOptions.Vars["integration_name"] = "Dockerhub Container Registry Updated"
		terraformOptions.Vars["non_os_package_support"] = true

		update := terraform.ApplyAndIdempotent(t, terraformOptions)
		updateData := GetContainerRegisteryDockerhub(update)
		assert.Equal(t, "Dockerhub Container Registry Updated", updateData.Data.Name)
		assert.Equal(t, true, updateData.Data.Data.NonOSPackageEval)
	}
}
