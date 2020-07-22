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
	testAccAlertChannelAwsCloudWatchResourceType = "lacework_alert_channel_aws_cloudwatch"
	testAccAlertChannelAwsCloudWatchResourceName = "example"

	// Environment variables for testing AWS CloudWatch Alert Channel Integrations
	testAccAlertChannelEventBusArn = "EVENT_BUS_ARN"
)

func TestAccAlertChannelAwsCloudWatch(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelAwsCloudWatchResourceType,
		testAccAlertChannelAwsCloudWatchResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelAwsCloudWatchEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelAwsCloudWatchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelAwsCloudWatchConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelAwsCloudWatchExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelAwsCloudWatchConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelAwsCloudWatchExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelAwsCloudWatchDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelAwsCloudWatchResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAwsCloudWatchAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.AwsCloudWatchIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelAwsCloudWatchExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetAwsCloudWatchAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.AwsCloudWatchIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.AwsCloudWatchIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelAwsCloudWatchEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelEventBusArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelEventBusArn)
	}
}

func testAccAlertChannelAwsCloudWatchConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    event_bus_arn = "%s"
}
`,
		testAccAlertChannelAwsCloudWatchResourceType,
		testAccAlertChannelAwsCloudWatchResourceName,
		enabled,
		os.Getenv(testAccAlertChannelEventBusArn),
	)
}
