package lacework

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/lacework/go-sdk/api"
)

const (
	testAccIntegrationOciCfgResourceType = "lacework_integration_oci_cfg"
	testAccIntegrationOciCfgResourceName = "example"

	// Environment variables for testing OCI Integrations
	testAccIntegrationOciEnvFingerprint = "OCI_FINGERPRINT"
	testAccIntegrationOciEnvPrivateKey  = "OCI_PRIVATE_KEY"
	testAccIntegrationOciEnvHomeRegion  = "OCI_HOME_REGION"
	testAccIntegrationOciEnvTenantID    = "OCI_TENANT_ID"
	testAccIntegrationOciEnvTenantName  = "OCI_TENANT_NAME"
	testAccIntegrationOciEnvUserOCID    = "OCI_USER_OCID"
)

func TestAccIntegrationOciCfg(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationOciCfgResourceType,
		testAccIntegrationOciCfgResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationOciCfgEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationOciCfgCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationOciCfgConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationOciCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationOciCfgConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationOciCfgExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationOciCfgCheckDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationOciCfgResourceType {
			continue
		}

		response, err := lacework.V2.CloudAccounts.GetOciCfg(rs.Primary.ID)
		if err != nil {
			if !strings.Contains(err.Error(), "[404] Not found") {
				return err
			}
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return fmt.Errorf("the OCI integration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIntegrationOciCfgExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.V2.CloudAccounts.GetOciCfg(rs.Primary.ID)
		if err != nil {
			return err
		}

		if response.Data.Name == "" {
			return fmt.Errorf("the OCI integration (%s) doesn't exist", rs.Primary.ID)
		}

		if response.Data.IntgGuid == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("the OCI integration (%s) doesn't exist", rs.Primary.ID)
	}
}

func testAccIntegrationOciCfgEnvVarsPreCheck(t *testing.T) {
	for _, envVar := range []string{
		testAccIntegrationOciEnvFingerprint,
		testAccIntegrationOciEnvPrivateKey,
		testAccIntegrationOciEnvHomeRegion,
		testAccIntegrationOciEnvTenantID,
		testAccIntegrationOciEnvTenantName,
		testAccIntegrationOciEnvUserOCID,
	} {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
}

func testAccIntegrationOciCfgConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    credentials {
        fingerprint = "%s"
        private_key = "%s"
    }
	home_region = "%s"
	tenant_id = "%s"
	tenant_name = "%s"
	user_ocid = "%s"
}
`,
		testAccIntegrationOciCfgResourceType,
		testAccIntegrationOciCfgResourceName,
		enabled,
		os.Getenv(testAccIntegrationOciEnvFingerprint),
		os.Getenv(testAccIntegrationOciEnvPrivateKey),
		os.Getenv(testAccIntegrationOciEnvHomeRegion),
		os.Getenv(testAccIntegrationOciEnvTenantID),
		os.Getenv(testAccIntegrationOciEnvTenantName),
		os.Getenv(testAccIntegrationOciEnvUserOCID),
	)
}
