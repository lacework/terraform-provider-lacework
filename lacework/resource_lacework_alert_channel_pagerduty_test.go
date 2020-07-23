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
	testAccAlertChannelPagerDutyResourceType = "lacework_alert_channel_pagerduty"
	testAccAlertChannelPagerDutyResourceName = "example"

	// Environment variables for testing PagerDuty Alert Channel Integrations
	testAccAlertChannelPagerDutyEnvURL = "INTEGRATION_KEY"
)

func TestAccAlertChannelPagerDuty(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelPagerDutyResourceType,
		testAccAlertChannelPagerDutyResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelPagerDutyEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelPagerDutyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelPagerDutyConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelPagerDutyExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelPagerDutyConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelPagerDutyExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelPagerDutyDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelPagerDutyResourceType {
			continue
		}

		response, err := lacework.Integrations.GetPagerDutyAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.PagerDutyIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelPagerDutyExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetPagerDutyAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.PagerDutyIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.PagerDutyIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelPagerDutyEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelPagerDutyEnvURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelPagerDutyEnvURL)
	}
}

func testAccAlertChannelPagerDutyConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    integration_key = "%s"
}
`,
		testAccAlertChannelPagerDutyResourceType,
		testAccAlertChannelPagerDutyResourceName,
		enabled,
		os.Getenv(testAccAlertChannelPagerDutyEnvURL),
	)
}
