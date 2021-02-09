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
	testAccAlertChannelMicrosoftTeamsResourceType = "lacework_alert_channel_microsoft_teams"
	testAccAlertChannelMicrosoftTeamsResourceName = "example"

	// Environment variables for testing Microsoft Teams Alert Channel Integrations
	testAccAlertChannelMicrosoftTeamsURL = "MICROSOFT_TEAMS_URL"
)

func TestAccAlertChannelMicrosoftTeams(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelMicrosoftTeamsResourceType,
		testAccAlertChannelMicrosoftTeamsResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelMicrosoftTeamsEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelMicrosoftTeamsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelMicrosoftTeamsConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelMicrosoftTeamsExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelMicrosoftTeamsConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelMicrosoftTeamsExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelMicrosoftTeamsDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelMicrosoftTeamsResourceType {
			continue
		}

		response, err := lacework.Integrations.GetMicrosoftTeamsAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.MicrosoftTeamsChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelMicrosoftTeamsExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetMicrosoftTeamsAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.MicrosoftTeamsChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.MicrosoftTeamsChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelMicrosoftTeamsEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelMicrosoftTeamsURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelMicrosoftTeamsURL)
	}
}

func testAccAlertChannelMicrosoftTeamsConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "integration test"
  enabled = %t
  teams_url = "%s"
}
		`,
		testAccAlertChannelMicrosoftTeamsResourceType,
		testAccAlertChannelMicrosoftTeamsResourceName,
		enabled,
		os.Getenv(testAccAlertChannelMicrosoftTeamsURL),
	)
}
