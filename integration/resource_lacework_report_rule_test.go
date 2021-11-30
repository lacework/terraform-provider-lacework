package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestReportRuleCreate applies integration terraform:
// => '../examples/resource_lacework_report_rule'
//
// It uses the go-sdk to verify the created report rule,
// applies an update and destroys it
func TestReportRuleCreate(t *testing.T) {
	name := fmt.Sprintf("Report Rule - %s", time.Now())
	resourceGroupName := fmt.Sprintf("Used for Report Rule Test - %s", time.Now())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_report_rule",
		Vars: map[string]interface{}{
			"name":                name,
			"description":         "Report Rule created by Terraform",
			"severities":          []string{"Critical"},
			"aws_pci":             true,
			"resource_group_name": resourceGroupName,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Report Rule
	terraformOptions.TimeBetweenRetries = 2 * time.Second
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetReportRuleProps(create)

	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualChannels := terraform.Output(t, terraformOptions, "channels")
	actualSeverities := terraform.Output(t, terraformOptions, "severities")
	actualAwsNotifications := terraform.Output(t, terraformOptions, "aws_pci")

	assert.Equal(t, "Report Rule created by Terraform", createProps.Data.Filter.Description)
	assert.Equal(t, []string{"Critical"}, api.NewReportRuleSeveritiesFromIntSlice(createProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, actualAwsNotifications, "true")
	assert.True(t, createProps.Data.ReportNotificationTypes.AwsPci)
	assert.Equal(t, "Report Rule created by Terraform", actualDescription)
	assert.Equal(t, string("[Critical]"), actualSeverities)
	assert.NotEmpty(t, actualChannels)
	assert.NotEmpty(t, createProps.Data.EmailAlertChannels)

	// Update Report Rule
	terraformOptions.Vars = map[string]interface{}{
		"name":                name,
		"description":         "Updated Report Rule created by Terraform",
		"severities":          []string{"High", "Medium"},
		"aws_pci":             false,
		"resource_group_name": resourceGroupName,
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetReportRuleProps(update)
	actualDescription = terraform.Output(t, terraformOptions, "description")
	actualChannels = terraform.Output(t, terraformOptions, "channels")
	actualSeverities = terraform.Output(t, terraformOptions, "severities")
	actualAwsNotifications = terraform.Output(t, terraformOptions, "aws_pci")

	assert.Equal(t, "Updated Report Rule created by Terraform", updateProps.Data.Filter.Description)
	assert.Equal(t, []string{"High", "Medium"}, api.NewReportRuleSeveritiesFromIntSlice(updateProps.Data.Filter.Severity).ToStringSlice())
	assert.Equal(t, "Updated Report Rule created by Terraform", actualDescription)
	assert.Equal(t, "[High Medium]", actualSeverities)
	assert.Equal(t, actualAwsNotifications, "false")
	assert.False(t, updateProps.Data.ReportNotificationTypes.AwsPci)
	assert.NotEmpty(t, actualChannels)
	assert.NotEmpty(t, updateProps.Data.EmailAlertChannels)
}
