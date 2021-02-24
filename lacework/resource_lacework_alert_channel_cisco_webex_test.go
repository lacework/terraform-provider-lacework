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
	testAccAlertChannelCiscoWebexResourceType = "lacework_alert_channel_cisco_webex"
	testAccAlertChannelCiscoWebexResourceName = "example"

	testAccAlertChannelCiscoWebexWebhookURL = "WEBHOOK"
)

func TestAccAlertChannelCiscoWebex(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelCiscoWebexResourceType,
		testAccAlertChannelCiscoWebexResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelCiscoWebexEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelCiscoWebexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelCiscoWebexConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelCiscoWebexExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelCiscoWebexConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelCiscoWebexExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelCiscoWebexDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelCiscoWebexResourceType {
			continue
		}

		response, err := lacework.Integrations.GetCiscoWebexAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.CiscoWebexChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelCiscoWebexExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetCiscoWebexAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.CiscoWebexChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.CiscoWebexChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelCiscoWebexEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelCiscoWebexWebhookURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelCiscoWebexWebhookURL)
	}
}

func testAccAlertChannelCiscoWebexConfig(enabled bool) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "integration test"
		enabled = %t
		webhook_url = "%s"
	}
	`,
		testAccAlertChannelCiscoWebexResourceType,
		testAccAlertChannelCiscoWebexResourceName,
		enabled,
		os.Getenv(testAccAlertChannelCiscoWebexWebhookURL),
	)
}
