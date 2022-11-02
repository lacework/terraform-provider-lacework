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
	testAccIntegrationAwsEksAuditResourceType = "lacework_integration_aws_eks_audit_log"
	testAccIntegrationAwsEksAuditResourceName = "example"

	// Environment variables for testing AWS_EKS_AUDIT only
	testAccIntegrationAwsEnvSnsArn = "sns:arn"
)

func TestAccIntegrationAwsEksAudit(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAwsEksAuditResourceType,
		testAccIntegrationAwsEksAuditResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAwsEksAuditEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAwsEksAuditDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAwsEksAuditConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsEksAuditExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAwsEksAuditConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAwsEksAuditExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAwsEksAuditDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAwsEksAuditResourceType {
			continue
		}

		response, err := lacework.V2.CloudAccounts.GetAwsEksAudit(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return fmt.Errorf("the AWS integration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIntegrationAwsEksAuditExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.V2.CloudAccounts.GetAwsEksAudit(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.Name == "" {
			return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("the AWS integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationAwsEksAuditEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAwsEnvSnsArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvSnsArn)
	}
	if v := os.Getenv(testAccIntegrationAwsEnvRoleArn); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvRoleArn)
	}
	if v := os.Getenv(testAccIntegrationAwsEnvExternalId); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAwsEnvExternalId)
	}
}

func testAccIntegrationAwsEksAuditConfig(enabled bool) string {
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
		testAccIntegrationAwsEksAuditResourceType,
		testAccIntegrationAwsEksAuditResourceName,
		enabled,
		os.Getenv(testAccIntegrationAwsEnvSnsArn),
		os.Getenv(testAccIntegrationAwsEnvRoleArn),
		os.Getenv(testAccIntegrationAwsEnvExternalId),
	)
}
