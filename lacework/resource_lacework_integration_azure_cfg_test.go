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
	testAccIntegrationAzureCfgResourceType = "lacework_integration_azure_cfg"
	testAccIntegrationAzureCfgResourceName = "example"

	// Environment variables for testing AZURE_CFG
	testAccIntegrationAzureCfgEnvTenantID     = "AZURE_TENANT_ID"
	testAccIntegrationAzureCfgEnvClientID     = "AZURE_CLIENT_ID"
	testAccIntegrationAzureCfgEnvClientSecret = "AZURE_CLIENT_SECRET"
)

func TestAccIntegrationAzureCfg(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAzureCfgResourceType,
		testAccIntegrationAzureCfgResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAzureCfgEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAzureCfgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAzureCfgConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAzureCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAzureCfgConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAzureCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAzureCfgDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAzureCfgResourceType {
			continue
		}

		response, err := lacework.Integrations.GetAzure(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf("the AZURE integration (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationAzureCfgExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetAzure(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the AZURE integration (%s) doesn't exist", rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the AZURE integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationAzureCfgEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAzureCfgEnvTenantID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureCfgEnvTenantID)
	}
	if v := os.Getenv(testAccIntegrationAzureCfgEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureCfgEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationAzureCfgEnvClientSecret); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureCfgEnvClientSecret)
	}
}

func testAccIntegrationAzureCfgConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    tenant_id = "%s"
    credentials {
        client_id = "%s"
        client_secret = "%s"
    }
}
`,
		testAccIntegrationAzureCfgResourceType,
		testAccIntegrationAzureCfgResourceName,
		enabled,
		os.Getenv(testAccIntegrationAzureCfgEnvTenantID),
		os.Getenv(testAccIntegrationAzureCfgEnvClientID),
		os.Getenv(testAccIntegrationAzureCfgEnvClientSecret),
	)
}
