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
	testAccAlertChannelJiraServerResourceType = "lacework_alert_channel_jira_server"
	testAccAlertChannelJiraServerResourceName = "example"

	// Environment variables for testing Jira ServerAlert Channel Integrations
	testAccAlertChannelJiraPassword = "JIRA_PASSWORD"
)

func TestAccAlertChannelJiraServer(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelJiraServerResourceType,
		testAccAlertChannelJiraServerResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelJiraServerEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelJiraServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelJiraServerConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelJiraServerExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelJiraServerConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelJiraServerExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelJiraServerDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelJiraServerResourceType {
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

func testAccCheckAlertChannelJiraServerExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccAlertChannelJiraServerEnvVarsPreCheck(t *testing.T) {
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
	if v := os.Getenv(testAccAlertChannelJiraPassword); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelJiraPassword)
	}
}

func testAccAlertChannelJiraServerConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name        = "integration test"
    enabled     = %t
    jira_url    = "%s"
    issue_type  = "%s"
    project_key = "%s"
    username    = "%s"
    password    = "%s"
    group_issues_by = "%s"
    custom_template_file = "%s"
}
`,
		testAccAlertChannelJiraServerResourceType,
		testAccAlertChannelJiraServerResourceName,
		enabled,
		os.Getenv(testAccAlertChannelJiraURL),
		os.Getenv(testAccAlertChannelJiraIssueType),
		os.Getenv(testAccAlertChannelJiraProjectKey),
		os.Getenv(testAccAlertChannelJiraUsername),
		os.Getenv(testAccAlertChannelJiraPassword),
		os.Getenv(testAccAlertChannelJiraGroupIssuesBy),
		os.Getenv(testAccAlertChannelJiraCustomTemplateFile),
	)
}
