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
	testAccAlertChannelSplunkResourceType = "lacework_alert_channel_splunk"
	testAccAlertChannelSplunkResourceName = "example"

	testAccAlertChannelSplunkChannel  = "SPLUNK_URL"
	testAccAlertChannelSplunkHecToken = "HEC_TOKEN"
	testAccAlertChannelSplunkHost     = "HOST"
	testAccAlertChannelSplunkPort     = "PORT"
	testAccAlertChannelSplunkSsl      = "SSL"
	testAccAlertChannelSplunkSource   = "SOURCE"
	testAccAlertChannelSplunkIndex    = "INDEX"
)

func TestAccAlertChannelSplunk(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelSplunkResourceType,
		testAccAlertChannelSplunkResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelSplunkEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelSplunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelSplunkConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelSplunkExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelSplunkConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelSplunkExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelSplunkDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelSplunkResourceType {
			continue
		}

		response, err := lacework.Integrations.GetSplunkAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.SplunkIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelSplunkExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetSplunkAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.SplunkIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.SplunkIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelSplunkEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelSplunkChannel); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkChannel)
	}
	if v := os.Getenv(testAccAlertChannelSplunkHecToken); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkHecToken)
	}
	if v := os.Getenv(testAccAlertChannelSplunkHost); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkHost)
	}
	if v := os.Getenv(testAccAlertChannelSplunkPort); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkPort)
	}
	if v := os.Getenv(testAccAlertChannelSplunkSsl); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkSsl)
	}
	if v := os.Getenv(testAccAlertChannelSplunkSource); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkSource)
	}
	if v := os.Getenv(testAccAlertChannelSplunkIndex); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelSplunkIndex)
	}
}

func testAccAlertChannelSplunkConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "integration test"
  enabled = %t
  channel = "%s"
  hec_token = "%s"
  host = "%s"
  port = "%s"
  ssl = "%s"
  event_data {
    index = "%s"
    source = "%s"
  }
}
	`,
		testAccAlertChannelSplunkResourceType,
		testAccAlertChannelSplunkResourceName,
		enabled,
		os.Getenv(testAccAlertChannelSplunkChannel),
		os.Getenv(testAccAlertChannelSplunkHecToken),
		os.Getenv(testAccAlertChannelSplunkHost),
		os.Getenv(testAccAlertChannelSplunkPort),
		os.Getenv(testAccAlertChannelSplunkSsl),
		os.Getenv(testAccAlertChannelSplunkIndex),
		os.Getenv(testAccAlertChannelSplunkSource),
	)
}
