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
	testAccIntegrationGCPCFGResourceType = "lacework_integration_gcp_cfg"
	testAccIntegrationGCPCFGResourceName = "example"

	// Environment variables for testing GCP CFG
	testAccIntegrationGCPCFGEnvClientID      = "GCP_CLIENT_ID"
	testAccIntegrationGCPCFGEnvPrivateKeyID  = "GCP_PRIVATE_KEY_ID"
	testAccIntegrationGCPCFGEnvPrivateKey    = "GCP_PRIVATE_KEY"
	testAccIntegrationGCPCFGEnvClientEmail   = "GCP_CLIENT_EMAIL"
	testAccIntegrationGCPCFGEnvResourceLevel = "GCP_RESOURCE_LEVEL"
	testAccIntegrationGCPCFGEnvResourceID    = "GCP_RESOURCE_ID"
)

func TestAccIntegrationGCPCFG(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationGCPCFGResourceType,
		testAccIntegrationGCPCFGResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationGCPCFGEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationGCPCFGDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationGCPCFGConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGCPCFGExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationGCPCFGConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGCPCFGExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationGCPCFGDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != testAccIntegrationGCPCFGResourceType {
			continue
		}

		response, err := lacework.GetGCPConfigIntegration(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the Google Cloud Platform integration (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationGCPCFGExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.GetGCPConfigIntegration(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the Google Cloud Platform integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the Google Cloud Platform integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationGCPCFGEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationGCPCFGEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGCPCFGEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationGCPCFGEnvPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGCPCFGEnvPrivateKeyID)
	}
	if v := os.Getenv(testAccIntegrationGCPCFGEnvPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGCPCFGEnvPrivateKey)
	}
	if v := os.Getenv(testAccIntegrationGCPCFGEnvClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGCPCFGEnvClientEmail)
	}
	if v := os.Getenv(testAccIntegrationGCPCFGEnvResourceID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGCPCFGEnvResourceID)
	}
}

func testAccIntegrationGCPCFGConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "Example-GCP-Integration"
    enabled = %t
    credentials {
        client_id = "%s"
        private_key_id = "%s"
        client_email = "%s"
        private_key = "%s"
    }
    resource_id = "%s"
  %s
}
`,
		testAccIntegrationGCPCFGResourceType,
		testAccIntegrationGCPCFGResourceName,
		enabled,
		os.Getenv(testAccIntegrationGCPCFGEnvClientID),
		os.Getenv(testAccIntegrationGCPCFGEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGCPCFGEnvClientEmail),
		os.Getenv(testAccIntegrationGCPCFGEnvPrivateKey),
		os.Getenv(testAccIntegrationGCPCFGEnvResourceID),
		resourceLevelAttrOrEmpty(
			os.Getenv(testAccIntegrationGCPCFGEnvResourceLevel),
		),
	)
}
func resourceLevelAttrOrEmpty(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("resource_level = \"%s\"", s)
}
