package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/lacework/go-sdk/lwdomain"
	"github.com/stretchr/testify/assert"
)

// TestTeamMemberStandalone applies integration terraform:
// => '../examples/resource_lacework_team_member_standalone'
//
// It uses the go-sdk to verify the created team member,
// applies an update with new description and destroys it
func TestTeamMemberStandalone(t *testing.T) {
	email := fmt.Sprintf("vatasha.white+%d@lacework.net", time.Now().Unix())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_team_member_standalone",
		Vars:         map[string]interface{}{"email": email},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Standalone Team Member
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	tm := GetTeamMember(create)
	assert.Equal(t, email, tm.UserName)
	assert.Equal(t, "Marvel Comics", tm.Props.Company)
	assert.Equal(t, "Shuri", tm.Props.FirstName)
	assert.Equal(t, "White", tm.Props.LastName)
	assert.False(t, tm.Props.AccountAdmin)

	// Update Standalone Team Member
	terraformOptions.Vars["first_name"] = "Vatasha"
	terraformOptions.Vars["administrator"] = true

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	tmUpdate := GetTeamMember(update)
	assert.Equal(t, email, tmUpdate.UserName)
	assert.Equal(t, "Marvel Comics", tm.Props.Company)
	assert.Equal(t, "Vatasha", tmUpdate.Props.FirstName)
	assert.Equal(t, "White", tmUpdate.Props.LastName)
	assert.True(t, tmUpdate.Props.AccountAdmin)

}

// TestTeamMemberOrg applies integration terraform:
// => '../examples/resource_lacework_team_member_organization'
//
// It uses the go-sdk to verify the created team member,
// applies an update with new description and destroys it
func TestTeamMemberOrg(t *testing.T) {
	if os.Getenv("CI_STANDALONE_ACCOUNT") != "" {
		t.Skip("skipping organizational account test")
	}
	account := os.Getenv("LW_ACCOUNT")
	email := fmt.Sprintf("vatasha.white+%d@lacework.net", time.Now().Unix())

	if domain, err := lwdomain.New(account); err == nil {
		account = domain.Account
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_team_member_organization",
		Vars: map[string]interface{}{
			"email":         email,
			"user_accounts": []string{account},
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new Org Team Member
	create := terraform.InitAndApply(t, terraformOptions)
	tm := GetOrgTeamMember(create)
	assert.Equal(t, email, tm.UserName)
	assert.Equal(t, "Pokemon International Company", tm.Props.Company)
	assert.Equal(t, "Vatasha", tm.Props.FirstName)
	assert.Equal(t, "White", tm.Props.LastName)

	// The second apply should be idempotent. Why?
	// Because the APIs doesn't return some fields
	terraform.ApplyAndIdempotent(t, terraformOptions)

	// Update Org Team Member
	terraformOptions.Vars["first_name"] = "Shuri"
	terraformOptions.Vars["user_accounts"] = []string{}
	terraformOptions.Vars["admin_accounts"] = []string{account}

	update := terraform.ApplyAndIdempotent(t, terraformOptions)
	tmUpdate := GetOrgTeamMember(update)
	assert.Equal(t, email, tmUpdate.UserName)
	assert.Equal(t, "Pokemon International Company", tm.Props.Company)
	assert.Equal(t, "Shuri", tmUpdate.Props.FirstName)
	assert.Equal(t, "White", tmUpdate.Props.LastName)
	// TODO check with search for list of admin accounts
}
