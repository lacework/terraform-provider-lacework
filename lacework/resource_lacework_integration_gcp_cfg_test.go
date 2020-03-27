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
	testAccIntegrationGcpCfgResourceType = "lacework_integration_gcp_cfg"
	testAccIntegrationGcpCfgResourceName = "example"

	// Environment variables for testing GCP CFG
	testAccIntegrationGcpCfgEnvClientID      = "GCP_CLIENT_ID"
	testAccIntegrationGcpCfgEnvPrivateKeyID  = "GCP_PRIVATE_KEY_ID"
	testAccIntegrationGcpCfgEnvPrivateKey    = "GCP_PRIVATE_KEY"
	testAccIntegrationGcpCfgEnvClientEmail   = "GCP_CLIENT_EMAIL"
	testAccIntegrationGcpCfgEnvResourceLevel = "GCP_RESOURCE_LEVEL"
	testAccIntegrationGcpCfgEnvResourceID    = "GCP_RESOURCE_ID"
)

func TestAccIntegrationGcpCfg(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationGcpCfgResourceType,
		testAccIntegrationGcpCfgResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationGcpCfgEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationGcpCfgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationGcpCfgConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationGcpCfgConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationGcpCfgDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationGcpCfgResourceType {
			continue
		}

		response, err := lacework.Integrations.GetGcp(rs.Primary.ID)
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

func testAccCheckIntegrationGcpCfgExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationGcpCfgEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationGcpCfgEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpCfgEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationGcpCfgEnvPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpCfgEnvPrivateKeyID)
	}
	if v := os.Getenv(testAccIntegrationGcpCfgEnvPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpCfgEnvPrivateKey)
	}
	if v := os.Getenv(testAccIntegrationGcpCfgEnvClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpCfgEnvClientEmail)
	}
	if v := os.Getenv(testAccIntegrationGcpCfgEnvResourceID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpCfgEnvResourceID)
	}
}

func testAccIntegrationGcpCfgConfig(enabled bool) string {
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
		testAccIntegrationGcpCfgResourceType,
		testAccIntegrationGcpCfgResourceName,
		enabled,
		os.Getenv(testAccIntegrationGcpCfgEnvClientID),
		os.Getenv(testAccIntegrationGcpCfgEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGcpCfgEnvClientEmail),
		os.Getenv(testAccIntegrationGcpCfgEnvPrivateKey),
		os.Getenv(testAccIntegrationGcpCfgEnvResourceID),
		resourceLevelAttrOrEmpty(
			os.Getenv(testAccIntegrationGcpCfgEnvResourceLevel),
		),
	)
}
func resourceLevelAttrOrEmpty(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("resource_level = \"%s\"", s)
}
