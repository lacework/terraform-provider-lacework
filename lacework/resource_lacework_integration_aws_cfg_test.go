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
	testAccIntegrationAwsCfgResourceType = "lacework_integration_aws_cfg"
	testAccIntegrationAwsCfgResourceName = "example"

	// Environment variables for testing AWS Integrations
	testAccIntegrationAwsEnvRoleArn    = "AWS_ROLE_ARN"
	testAccIntegrationAwsEnvExternalId = "AWS_EXTERNAL_ID"
)

func TestAccIntegrationAwsCfg(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAwsCfgResourceType,
		testAccIntegrationAwsCfgResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAwsCfgEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAwsCfgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAwsCfgConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAwsCfgConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAwsCfgDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAwsCfgResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAws(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the AWS integration (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationAwsCfgExists(resourceTypeAndName string) resource.TestCheckFunc {
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
			return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationAwsCfgEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAwsEnvRoleArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvRoleArn)
	}
	if v := os.Getenv(testAccIntegrationAwsEnvExternalId); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvExternalId)
	}
}

func testAccIntegrationAwsCfgConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    credentials {
        role_arn = "%s"
        external_id = "%s"
    }
}
`,
		testAccIntegrationAwsCfgResourceType,
		testAccIntegrationAwsCfgResourceName,
		enabled,
		os.Getenv(testAccIntegrationAwsEnvRoleArn),
		os.Getenv(testAccIntegrationAwsEnvExternalId),
	)
}
