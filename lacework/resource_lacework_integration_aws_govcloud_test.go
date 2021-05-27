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
	testAccIntegrationAwsGovCloudCfgResourceType = "lacework_integration_aws_govcloud_cfg"
	testAccIntegrationAwsGovCloudCfgResourceName = "example"

	// Environment variables for testing AWS Integrations
	testAccIntegrationAwsGovCloudCfgAccessKeyID     = "AWS_ACCESS_KEY"
	testAccIntegrationAwsGovCloudCfgSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	testAccIntegrationAwsGovCloudCfgAccountID       = "AWS_ACCOUNT_ID"
)

func TestAccIntegrationAwsGovCloudCfg(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAwsGovCloudCfgResourceType,
		testAccIntegrationAwsGovCloudCfgResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAwsGovCloudCfgEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAwsGovCloudCfgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAwsGovCloudCfgConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsGovCloudCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAwsGovCloudCfgConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsGovCloudCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAwsGovCloudCfgDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAwsGovCloudCfgResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAws(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the AWS Gov Cloud integration (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationAwsGovCloudCfgExists(resourceTypeAndName string) resource.TestCheckFunc {
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
			return fmt.Errorf("the AWS Gov Cloud integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationAwsGovCloudCfgEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAwsGovCloudCfgAccessKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCloudCfgAccessKeyID)
	}
	if v := os.Getenv(testAccIntegrationAwsGovCloudCfgSecretAccessKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCloudCfgSecretAccessKey)
	}

	if v := os.Getenv(testAccIntegrationAwsGovCloudCfgAccountID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsGovCloudCfgAccountID)
	}
}

func testAccIntegrationAwsGovCloudCfgConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    credentials {
        access_key_id = "%s"
        secret_access_key = "%s"
    }
	account_id = "%s"
}
`,
		testAccIntegrationAwsGovCloudCfgResourceType,
		testAccIntegrationAwsGovCloudCfgResourceName,
		enabled,
		os.Getenv(testAccIntegrationAwsGovCloudCfgAccessKeyID),
		os.Getenv(testAccIntegrationAwsGovCloudCfgSecretAccessKey),
		os.Getenv(testAccIntegrationAwsGovCloudCfgAccountID),
	)
}
