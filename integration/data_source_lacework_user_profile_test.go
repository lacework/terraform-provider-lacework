package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestUserProfileDataSource uses the Terraform plan at:
// => '../examples/data_source_lacework_user_profile'
func TestUserProfileDataSource(t *testing.T) {
	// Read API access token
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/data_source_lacework_user_profile",
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	url := terraform.Output(t, terraformOptions, "lacework_user_profile_url")
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "lacework.net")
}
