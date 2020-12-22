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
	testAccAlertChannelAwsS3ResourceType = "lacework_alert_channel_aws_s3"
	testAccAlertChannelAwsS3ResourceName = "example"

	// Environment variables for testing Aws S3 Alert Channel Integrations
	testAccAlertChannelAwsS3ExternalID = "AWS_S3_EXTERNAL_ID"
	testAccAlertChannelAwsS3BucketArn  = "AWS_S3_BUCKET_ARN"
	testAccAlertChannelAwsS3RoleArn    = "AWS_S3_ROLE_ARN"
)

func TestAccAlertChannelAwsS3(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelAwsS3ResourceType,
		testAccAlertChannelAwsS3ResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelAwsS3EnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelAwsS3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelAwsS3Config(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelAwsS3Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelAwsS3Config(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelAwsS3Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelAwsS3Destroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelAwsS3ResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAwsS3AlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.AwsS3ChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelAwsS3Exists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetAwsS3AlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.AwsS3ChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.AwsS3ChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelAwsS3EnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelAwsS3ExternalID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelAwsS3ExternalID)
	}
	if v := os.Getenv(testAccAlertChannelAwsS3RoleArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelAwsS3RoleArn)
	}
	if v := os.Getenv(testAccAlertChannelAwsS3BucketArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelAwsS3BucketArn)
	}
}

func testAccAlertChannelAwsS3Config(enabled bool) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "integration test"
			enabled = %t
			bucket_arn = "%s"
			credentials {
			external_id = "%s"
			role_arn = "%s"
			}
		}
		`,
		testAccAlertChannelAwsS3ResourceType,
		testAccAlertChannelAwsS3ResourceName,
		enabled,
		os.Getenv(testAccAlertChannelAwsS3ExternalID),
		os.Getenv(testAccAlertChannelAwsS3RoleArn),
		os.Getenv(testAccAlertChannelAwsS3BucketArn),
	)
}
