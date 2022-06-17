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
	testAccIntegrationGcpGkeAuditResourceType = "GcpGkeAudit"
	testAccIntegrationGcpGkeAuditResourceName = "example"

	// Environment variables for testing GCP GKE Audit
	testAccIntegrationGcpGkeIntegrationType = "Project"
	testAccIntegrationGcpGkeProjectID       = "GCP_GKE_PROJECT"
	testAccIntegrationGcpGkeOrganizationID  = "GCP_GKE_ORG"
	testAccIntegrationGcpGkeSubscription    = "GCP_SUBSCRIPTION"
)

func TestAccIntegrationGcpGkeAudit(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationGcpGkeAuditResourceType,
		testAccIntegrationGcpGkeAuditResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationGcpGkeAuditEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationGcpGkeAuditDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationGcpGkeAuditConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpGkeAuditExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationGcpGkeAuditConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpGkeAuditExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationGcpGkeAuditDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != testAccIntegrationGcpGkeAuditResourceType {
			continue
		}

		response, err := lacework.Integrations.GetGcp(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the Google GKE Audit integration (%s) still exists",
					rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationGcpGkeAuditExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetGcp(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the Google GKE Audit integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the Google GKE Audit integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationGcpGkeAuditEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationGcpEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationGcpEnvPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvPrivateKeyID)
	}
	if v := os.Getenv(testAccIntegrationGcpEnvPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvPrivateKey)
	}
	if v := os.Getenv(testAccIntegrationGcpEnvClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvClientEmail)
	}
}

func testAccIntegrationGcpGkeAuditConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
   name = "Example-GCP-GKE-Audit-Integration"
   enabled = %t
   credentials {
       client_id = "%s"
       client_email = "%s"
       private_key_id = "%s"
       private_key = "%s"
   }
   integration_type = "%s"
   project_id = "%s"
   organization_id = "%s"
   subscription = "%s"
}
`,
		testAccIntegrationGcpGkeAuditResourceType,
		testAccIntegrationGcpGkeAuditResourceName,
		enabled,
		os.Getenv(testAccIntegrationGcpEnvClientID),
		os.Getenv(testAccIntegrationGcpEnvClientEmail),
		os.Getenv(testAccIntegrationGcpEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGcpEnvPrivateKey),
		testAccIntegrationGcpGkeIntegrationType,
		testAccIntegrationGcpGkeProjectID,
		testAccIntegrationGcpGkeOrganizationID,
		testAccIntegrationGcpGkeSubscription,
	)
}
