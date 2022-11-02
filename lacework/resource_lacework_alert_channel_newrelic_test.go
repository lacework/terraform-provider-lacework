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
	testAccAlertChannelNewRelicResourceType = "lacework_alert_channel_newrelic"
	testAccAlertChannelNewRelicResourceName = "example"

	testAccAlertChannelNewRelicAccountID = "ACCOUNT_ID"
	testAccAlertChannelNewRelicInsertKey = "INSERT_KEY"
)

func TestAccAlertChannelNewRelic(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelNewRelicResourceType,
		testAccAlertChannelNewRelicResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelNewRelicEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelNewRelicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelNewRelicConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelNewRelicExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelNewRelicConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelNewRelicExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelNewRelicDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelNewRelicResourceType {
			continue
		}

		response, err := lacework.V2.AlertChannels.GetNewRelicInsights(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return fmt.Errorf(
				"the %s integration (%s) still exists",
				api.NewRelicInsightsAlertChannelType, rs.Primary.ID,
			)
		}
	}

	return nil
}

func testAccCheckAlertChannelNewRelicExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.V2.AlertChannels.GetNewRelicInsights(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.Name == "" {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.NewRelicInsightsAlertChannelType, rs.Primary.ID)
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.NewRelicInsightsAlertChannelType, rs.Primary.ID)
	}
}

func testAccAlertChannelNewRelicEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelNewRelicAccountID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelNewRelicAccountID)
	}
	if v := os.Getenv(testAccAlertChannelNewRelicInsertKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelNewRelicInsertKey)
	}
}

func testAccAlertChannelNewRelicConfig(enabled bool) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "integration test"
		enabled = %t
		account_id = "%s"
		insert_key = "%s"
	}
	`,
		testAccAlertChannelNewRelicResourceType,
		testAccAlertChannelNewRelicResourceName,
		enabled,
		os.Getenv(testAccAlertChannelNewRelicAccountID),
		os.Getenv(testAccAlertChannelNewRelicInsertKey),
	)
}
