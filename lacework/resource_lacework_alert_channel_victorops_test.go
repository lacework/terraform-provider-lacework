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
	testAccAlertChannelVictorOpsResourceType = "lacework_alert_channel_victorops"
	testAccAlertChannelVictorOpsResourceName = "example"

	testAccAlertChannelVictorOpsWebhookURL = "INTG_URL"
)

func TestAccAlertChannelVictorOps(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelVictorOpsResourceType,
		testAccAlertChannelVictorOpsResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelVictorOpsEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelVictorOpsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelVictorOpsConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelVictorOpsExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelVictorOpsConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelVictorOpsExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelVictorOpsDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelVictorOpsResourceType {
			continue
		}

		response, err := lacework.Integrations.GetVictorOpsAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.VictorOpsChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelVictorOpsExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetVictorOpsAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.VictorOpsChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.VictorOpsChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelVictorOpsEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelVictorOpsWebhookURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelVictorOpsWebhookURL)
	}
}

func testAccAlertChannelVictorOpsConfig(enabled bool) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "integration test"
			enabled = %t
			webhook_url = "%s"
		}
		`,
		testAccAlertChannelVictorOpsResourceType,
		testAccAlertChannelVictorOpsResourceName,
		enabled,
		os.Getenv(testAccAlertChannelVictorOpsWebhookURL),
	)
}
