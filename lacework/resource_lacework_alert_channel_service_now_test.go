package lacework

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lacework/go-sdk/api"
)

const (
	testAccAlertChannelServiceNowResourceType = "lacework_alert_channel_splunk"
	testAccAlertChannelServiceNowResourceName = "example"

	testAccAlertChannelServiceNowChannel       = "SERVICE_NOW_REST"
	testAccAlertChannelServiceNowInstanceURL   = "INSTANCE_URL"
	testAccAlertChannelServiceNowUsername      = "USERNAME"
	testAccAlertChannelServiceNowPassword      = "PASSWORD"
	testAccAlertChannelServiceNowIssueGrouping = "ISSUE_GROUPING"
)

func TestAccAlertChannelServiceNow(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelServiceNowResourceType,
		testAccAlertChannelServiceNowResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelServiceNowEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelServiceNowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelServiceNowConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelServiceNowExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelServiceNowConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelServiceNowExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelServiceNowDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelServiceNowResourceType {
			continue
		}

		response, err := lacework.Integrations.GetServiceNowAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.ServiceNowChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelServiceNowExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetServiceNowAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.ServiceNowChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.ServiceNowChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelServiceNowEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelServiceNowChannel); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelServiceNowChannel)
	}
	if v := os.Getenv(testAccAlertChannelServiceNowInstanceURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelServiceNowInstanceURL)
	}
	if v := os.Getenv(testAccAlertChannelServiceNowUsername); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelServiceNowUsername)
	}
	if v := os.Getenv(testAccAlertChannelServiceNowPassword); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelServiceNowPassword)
	}
	if v := os.Getenv(testAccAlertChannelServiceNowIssueGrouping); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelServiceNowIssueGrouping)
	}
}

func testAccAlertChannelServiceNowConfig(enabled bool) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "integration test"
		enabled = %t
		channel = "%s"
		instance_url = "%s"
		username = "%s"
		password = "%s"
		issue_grouping = "%s"
	}
	`,
		testAccAlertChannelServiceNowResourceType,
		testAccAlertChannelServiceNowResourceName,
		enabled,
		os.Getenv(testAccAlertChannelServiceNowChannel),
		os.Getenv(testAccAlertChannelServiceNowInstanceURL),
		os.Getenv(testAccAlertChannelServiceNowUsername),
		os.Getenv(testAccAlertChannelServiceNowPassword),
		os.Getenv(testAccAlertChannelServiceNowIssueGrouping),
	)
}
