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
	testAccIntegrationAwsCloudTrailResourceType = "lacework_integration_aws_ct"
	testAccIntegrationAwsCloudTrailResourceName = "example"

	// Environment variables for testing AWS_CT_SQS only
	testAccIntegrationAwsEnvQueueUrl = "AWS_QUEUE_URL"
)

func TestAccIntegrationAwsCloudTrail(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAwsCloudTrailResourceType,
		testAccIntegrationAwsCloudTrailResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAwsCloudTrailEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAwsCloudTrailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAwsCloudTrailConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsCloudTrailExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAwsCloudTrailConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsCloudTrailExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAwsCloudTrailDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAwsCloudTrailResourceType {
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

func testAccCheckIntegrationAwsCloudTrailExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationAwsCloudTrailEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAwsEnvQueueUrl); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvQueueUrl)
	}
	if v := os.Getenv(testAccIntegrationAwsEnvRoleArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvRoleArn)
	}
	if v := os.Getenv(testAccIntegrationAwsEnvExternalId); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvExternalId)
	}
}

func testAccIntegrationAwsCloudTrailConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    queue_url = %s
    credentials {
        role_arn = "%s"
        external_id = "%s"
    }
}
`,
		testAccIntegrationAwsCloudTrailResourceType,
		testAccIntegrationAwsCloudTrailResourceName,
		enabled,
		os.Getenv(testAccIntegrationAwsEnvQueueUrl),
		os.Getenv(testAccIntegrationAwsEnvRoleArn),
		os.Getenv(testAccIntegrationAwsEnvExternalId),
	)
}
