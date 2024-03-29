package integration

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestPolicyCreate applies integration terraform:
// => '../examples/resource_lacework_policy'
//
// It uses the go-sdk to verify the created policy,
// applies an update and destroys it
// nolint
func TestPolicyCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_policy",
		EnvVars:      tokenEnvVar,
		Vars: map[string]interface{}{
			"title":       "lql-terraform-policy",
			"severity":    "High",
			"type":        "Violation",
			"description": "Policy Created via Terraform",
			"remediation": "Please Investigate",
			"evaluation":  "Hourly",
			"tags":        []string{"cloud_AWS", "resource_Cloudtrail"},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Policy
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetPolicyProps(create)

	actualTitle := terraform.Output(t, terraformOptions, "title")
	actualSeverity := terraform.Output(t, terraformOptions, "severity")
	actualType := terraform.Output(t, terraformOptions, "type")
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualRemediation := terraform.Output(t, terraformOptions, "remediation")
	actualEvaluation := terraform.Output(t, terraformOptions, "evaluation")
	actualTags := terraform.Output(t, terraformOptions, "tags")

	assert.Contains(t, "lql-terraform-policy", createProps.Data.Title)
	assert.Contains(t, "high", createProps.Data.Severity)
	assert.Contains(t, "Violation", createProps.Data.PolicyType)
	assert.Contains(t, "Policy Created via Terraform", createProps.Data.Description)
	assert.Contains(t, "Please Investigate", createProps.Data.Remediation)
	assert.Contains(t, "Hourly", createProps.Data.EvalFrequency)
	assert.NotContains(t, createProps.Data.Tags, "custom")
	assert.Contains(t, createProps.Data.Tags, "cloud_AWS")
	assert.Contains(t, createProps.Data.Tags, "resource_Cloudtrail")
	// @afiune Tags are now restricted and we can't use domain:XYZ or subdomain:ABC
	// assert.ElementsMatch(t, []string{"domain:AWS", "subdomain:Cloudtrail"}, createProps.Data.Tags)

	assert.Equal(t, "lql-terraform-policy", actualTitle)
	assert.Equal(t, "high", actualSeverity)
	assert.Equal(t, "Violation", actualType)
	assert.Equal(t, "Policy Created via Terraform", actualDescription)
	assert.Equal(t, "Please Investigate", actualRemediation)
	assert.Equal(t, "Hourly", actualEvaluation)
	assert.Equal(t, "[cloud_AWS resource_Cloudtrail]", actualTags)

	// Update Policy
	terraformOptions.Vars = map[string]interface{}{
		"title":       "lql-terraform-policy-updated",
		"severity":    "Low",
		"description": "Policy Created via Terraform Updated",
		"remediation": "Please Ignore",
		"evaluation":  "Daily",
		"tags":        []string{"cloud_AWS", "resource_Cloudtrail", "custom"},
	}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	updateProps := GetPolicyProps(update)

	actualTitle = terraform.Output(t, terraformOptions, "title")
	actualSeverity = terraform.Output(t, terraformOptions, "severity")
	actualType = terraform.Output(t, terraformOptions, "type")
	actualDescription = terraform.Output(t, terraformOptions, "description")
	actualRemediation = terraform.Output(t, terraformOptions, "remediation")
	actualEvaluation = terraform.Output(t, terraformOptions, "evaluation")
	actualTags = terraform.Output(t, terraformOptions, "tags")

	assert.Contains(t, "lql-terraform-policy-updated", updateProps.Data.Title)
	assert.Contains(t, "low", updateProps.Data.Severity)
	assert.Contains(t, "Policy Created via Terraform Updated", updateProps.Data.Description)
	assert.Contains(t, "Please Ignore", updateProps.Data.Remediation)
	assert.Contains(t, "Daily", updateProps.Data.EvalFrequency)
	assert.Contains(t, updateProps.Data.Tags, "custom")
	assert.Contains(t, updateProps.Data.Tags, "cloud_AWS")
	assert.Contains(t, updateProps.Data.Tags, "resource_Cloudtrail")
	// @afiune Tags are now restricted and we can't use domain:XYZ or subdomain:ABC
	// assert.ElementsMatch(t, []string{"custom", "domain:AWS", "subdomain:Cloudtrail"}, updateProps.Data.Tags)

	assert.Equal(t, "lql-terraform-policy-updated", actualTitle)
	assert.Equal(t, "low", actualSeverity)
	assert.Equal(t, "Violation", actualType)
	assert.Equal(t, "Policy Created via Terraform Updated", actualDescription)
	assert.Equal(t, "Please Ignore", actualRemediation)
	assert.Equal(t, "Daily", actualEvaluation)
	// @afiune this output is not consistent, CI is getting.
	//
	// Error: Not equal:
	//           expected: "[custom cloud_AWS resource_Cloudtrail]"
	//           actual  : "[cloud_AWS custom resource_Cloudtrail]"
	//
	// assert.Equal(t, "[custom cloud_AWS resource_Cloudtrail]", actualTags)
}

func TestPolicyCreateWithPolicyIDSuffix(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	suffix := fmt.Sprintf("terraform-%d", rand.Intn(1000))
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_policy",
		Vars: map[string]interface{}{
			"title":            "lql-terraform-policy",
			"policy_id_suffix": suffix,
			"severity":         "High",
			"type":             "Violation",
			"description":      "Policy Created via Terraform",
			"remediation":      "Please Investigate",
			"evaluation":       "Hourly",
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Policy
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	createProps := GetPolicyProps(create)

	actualTitle := terraform.Output(t, terraformOptions, "title")
	actualSeverity := terraform.Output(t, terraformOptions, "severity")
	actualType := terraform.Output(t, terraformOptions, "type")
	actualDescription := terraform.Output(t, terraformOptions, "description")
	actualRemediation := terraform.Output(t, terraformOptions, "remediation")
	actualEvaluation := terraform.Output(t, terraformOptions, "evaluation")
	actualSuffix := terraform.Output(t, terraformOptions, "policy_id_suffix")

	assert.Contains(t, "lql-terraform-policy", createProps.Data.Title)
	assert.Contains(t, "high", createProps.Data.Severity)
	assert.Contains(t, "Violation", createProps.Data.PolicyType)
	assert.Contains(t, "Policy Created via Terraform", createProps.Data.Description)
	assert.Contains(t, "Please Investigate", createProps.Data.Remediation)
	assert.Contains(t, "Hourly", createProps.Data.EvalFrequency)

	assert.Equal(t, "lql-terraform-policy", actualTitle)
	assert.Equal(t, "high", actualSeverity)
	assert.Equal(t, "Violation", actualType)
	assert.Equal(t, "Policy Created via Terraform", actualDescription)
	assert.Equal(t, "Please Investigate", actualRemediation)
	assert.Equal(t, "Hourly", actualEvaluation)
	assert.Contains(t, suffix, actualSuffix)

	// Update Policy
	terraformOptions.Vars = map[string]interface{}{
		"title":            "lql-terraform-policy-updated",
		"policy_id_suffix": "modified-id-suffix",
		"severity":         "Low",
		"description":      "Policy Created via Terraform Updated",
		"remediation":      "Please Ignore",
		"evaluation":       "Daily",
	}

	msg, err := terraform.ApplyE(t, terraformOptions)

	assert.Error(t, err)
	assert.Contains(t, msg, "unable to change ID of an existing policy")
}
