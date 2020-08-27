package lacework

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lacework/go-sdk/api"
)

const (
	testAccAlertChannelJiraCloudResourceType = "lacework_alert_channel_jira_cloud"
	testAccAlertChannelJiraCloudResourceName = "example"

	// Environment variables for testing Jira Alert Channel Integrations
	testAccAlertChannelJiraURL                = "JIRA_URL"
	testAccAlertChannelJiraIssueType          = "JIRA_ISSUE_TYPE"
	testAccAlertChannelJiraProjectKey         = "JIRA_PROJECT_KEY"
	testAccAlertChannelJiraUsername           = "JIRA_USERNAME"
	testAccAlertChannelJiraGroupIssuesBy      = "JIRA_GROUP_ISSUES_BY"
	testAccAlertChannelJiraCustomTemplateFile = "JIRA_CUSTOM_TEMPLATE_FILE"

	// this is only for Jira Cloud
	testAccAlertChannelJiraApiToken = "JIRA_API_TOKEN"
)

func TestAccAlertChannelJiraCloud(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelJiraCloudResourceType,
		testAccAlertChannelJiraCloudResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelJiraCloudEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelJiraCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelJiraCloudConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelJiraCloudExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelJiraCloudConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelJiraCloudExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelJiraCloudDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelJiraCloudResourceType {
			continue
		}

		response, err := lacework.Integrations.GetJiraAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.JiraIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelJiraCloudExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetJiraAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.JiraIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.JiraIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelJiraCloudEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelJiraURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraURL)
	}
	if v := os.Getenv(testAccAlertChannelJiraIssueType); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraIssueType)
	}
	if v := os.Getenv(testAccAlertChannelJiraProjectKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraProjectKey)
	}
	if v := os.Getenv(testAccAlertChannelJiraUsername); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraUsername)
	}
	if v := os.Getenv(testAccAlertChannelJiraGroupIssuesBy); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraGroupIssuesBy)
	}
	if v := os.Getenv(testAccAlertChannelJiraCustomTemplateFile); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraCustomTemplateFile)
	}
	if v := os.Getenv(testAccAlertChannelJiraApiToken); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraApiToken)
	}
}

func testAccAlertChannelJiraCloudConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name        = "integration test"
    enabled     = %t
    jira_url    = "%s"
    issue_type  = "%s"
    project_key = "%s"
    username    = "%s"
    api_token   = "%s"
    group_issues_by = "%s"
    custom_template_file = "%s"
}
`,
		testAccAlertChannelJiraCloudResourceType,
		testAccAlertChannelJiraCloudResourceName,
		enabled,
		os.Getenv(testAccAlertChannelJiraURL),
		os.Getenv(testAccAlertChannelJiraIssueType),
		os.Getenv(testAccAlertChannelJiraProjectKey),
		os.Getenv(testAccAlertChannelJiraUsername),
		os.Getenv(testAccAlertChannelJiraApiToken),
		os.Getenv(testAccAlertChannelJiraGroupIssuesBy),
		os.Getenv(testAccAlertChannelJiraCustomTemplateFile),
	)
}
