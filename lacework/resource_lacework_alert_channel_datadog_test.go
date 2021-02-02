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
	testAccAlertChannelDatadogResourceType = "lacework_alert_channel_datadog"
	testAccAlertChannelDatadogResourceName = "example"

	testAccAlertChannelDatadogSite    = "DATADOG_SITE"
	testAccAlertChannelDatadogService = "DATADOG_TYPE"
	testAccAlertChannelDatadogApiKey  = "API_KEY"
)

func TestAccAlertChannelDatadog(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelDatadogResourceType,
		testAccAlertChannelDatadogResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelDatadogEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelDatadogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelDatadogConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelDatadogExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelDatadogConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelDatadogExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelDatadogDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelDatadogResourceType {
			continue
		}

		response, err := lacework.Integrations.GetDatadogAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.DatadogIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelDatadogExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetDatadogAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.DatadogIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.DatadogIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelDatadogEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelDatadogSite); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelDatadogSite)
	}
	if v := os.Getenv(testAccAlertChannelDatadogService); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelDatadogService)
	}
	if v := os.Getenv(testAccAlertChannelDatadogApiKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelDatadogApiKey)
	}
}

func testAccAlertChannelDatadogConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "integration test"
  enabled = %t
  datadog_site = "%s"
  datadog_site = "%s"
  api_key = "%s"
}
	`,
		testAccAlertChannelDatadogResourceType,
		testAccAlertChannelDatadogResourceName,
		enabled,
		os.Getenv(testAccAlertChannelDatadogSite),
		os.Getenv(testAccAlertChannelDatadogService),
		os.Getenv(testAccAlertChannelDatadogApiKey),
	)
}
