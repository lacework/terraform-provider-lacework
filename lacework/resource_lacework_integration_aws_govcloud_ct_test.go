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
	testAccIntegrationAwsGovCloudCTResourceType = "lacework_integration_aws_govcloud_ct"
	testAccIntegrationAwsGovCloudCTResourceName = "example"

	// Environment variables for testing AWS Integrations
	testAccIntegrationAwsGovCTQueueUrl        = "AWS_QUEUE_URL"
	testAccIntegrationAwsGovCTAccessKeyID     = "AWS_ACCESS_KEY"
	testAccIntegrationAwsGovCTSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	testAccIntegrationAwsGovCTAccountID       = "AWS_ACCOUNT_ID"
)

func TestAccIntegrationAwsGovCloudCT(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAwsGovCloudCTResourceType,
		testAccIntegrationAwsGovCloudCTResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAwsGovCloudCTEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAwsGovCloudCTDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAwsGovCloudCTConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsGovCloudCTExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAwsGovCloudCTConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsGovCloudCTExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAwsGovCloudCTDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAwsGovCloudCTResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAws(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the AWS Gov Cloud CloudTrail integration (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationAwsGovCloudCTExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetAws(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the AWS Gov Cloud CloudTrail integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationAwsGovCloudCTEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAwsGovCTQueueUrl); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCTQueueUrl)
	}

	if v := os.Getenv(testAccIntegrationAwsGovCTAccessKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCTAccessKeyID)
	}
	if v := os.Getenv(testAccIntegrationAwsGovCTSecretAccessKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCTSecretAccessKey)
	}

	if v := os.Getenv(testAccIntegrationAwsGovCTAccountID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCTAccountID)
	}
}

func testAccIntegrationAwsGovCloudCTConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
	queue_url = "%s"
    credentials {
        access_key_id = "%s"
        secret_access_key = "%s"
        account_id = "%s"
    }
}
`,
		testAccIntegrationAwsGovCloudCTResourceType,
		testAccIntegrationAwsGovCloudCTResourceName,
		enabled,
		os.Getenv(testAccIntegrationAwsGovCTQueueUrl),
		os.Getenv(testAccIntegrationAwsGovCTAccessKeyID),
		os.Getenv(testAccIntegrationAwsGovCTSecretAccessKey),
		os.Getenv(testAccIntegrationAwsGovCTAccountID),
	)
}
