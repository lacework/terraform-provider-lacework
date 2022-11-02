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
	testAccAlertChannelQRadarResourceType = "lacework_alert_channel_qradar"
	testAccAlertChannelQRadarResourceName = "example"

	testAccAlertChannelQRadarCommunicationType = "QRADAR_COMM_TYPE"
	testAccAlertChannelQRadarHostPort          = "QRADAR_HOST_PORT"
	testAccAlertChannelQRadarHostURL           = "QRADAR_HOST_URL"
)

func TestAccAlertChannelQRadar(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelQRadarResourceType,
		testAccAlertChannelQRadarResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelQRadarEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelQRadarDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelQRadarConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelQRadarExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelQRadarConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelQRadarExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelQRadarDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelQRadarResourceType {
			continue
		}

		response, err := lacework.V2.AlertChannels.GetIbmQRadar(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return fmt.Errorf(
				"the %s integration (%s) still exists",
				api.QRadarChannelIntegration, rs.Primary.ID,
			)
		}
	}

	return nil
}

func testAccCheckAlertChannelQRadarExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.V2.AlertChannels.GetIbmQRadar(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.Name == "" {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.QRadarChannelIntegration, rs.Primary.ID)
		}

		if response.Data.ID() == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.QRadarChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelQRadarEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelQRadarHostURL); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelQRadarHostURL)
	}
	if v := os.Getenv(testAccAlertChannelQRadarHostPort); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelQRadarHostPort)
	}
	if v := os.Getenv(testAccAlertChannelQRadarCommunicationType); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelQRadarCommunicationType)
	}
}

func testAccAlertChannelQRadarConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "integration test"
  enabled = %t
  host_url = "%s"
  host_port = "%s"
  communication_type = "%s"
}
	`,
		testAccAlertChannelQRadarResourceType,
		testAccAlertChannelQRadarResourceName,
		enabled,
		os.Getenv(testAccAlertChannelQRadarHostURL),
		os.Getenv(testAccAlertChannelQRadarHostPort),
		os.Getenv(testAccAlertChannelQRadarCommunicationType),
	)
}
