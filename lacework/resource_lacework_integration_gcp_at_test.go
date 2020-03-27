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
	testAccIntegrationGcpAtResourceType = "lacework_integration_gcp_at"
	testAccIntegrationGcpAtResourceName = "example"

	// Environment variables for testing GCP CFG
	testAccIntegrationGcpAtEnvClientID      = "GCP_CLIENT_ID"
	testAccIntegrationGcpAtEnvPrivateKeyID  = "GCP_PRIVATE_KEY_ID"
	testAccIntegrationGcpAtEnvPrivateKey    = "GCP_PRIVATE_KEY"
	testAccIntegrationGcpAtEnvClientEmail   = "GCP_CLIENT_EMAIL"
	testAccIntegrationGcpAtEnvResourceLevel = "GCP_RESOURCE_LEVEL"
	testAccIntegrationGcpAtEnvResourceID    = "GCP_RESOURCE_ID"
	testAccIntegrationGcpAtEnvSubscription  = "GCP_SUBSCRIPTION"
)

func TestAccIntegrationGcpAt(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationGcpAtResourceType,
		testAccIntegrationGcpAtResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationGcpAtEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationGcpAtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationGcpAtConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpAtExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationGcpAtConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcpAtExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationGcpAtDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != testAccIntegrationGcpAtResourceType {
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

func testAccCheckIntegrationGcpAtExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationGcpAtEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationGcpAtEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationGcpAtEnvPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvPrivateKeyID)
	}
	if v := os.Getenv(testAccIntegrationGcpAtEnvPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvPrivateKey)
	}
	if v := os.Getenv(testAccIntegrationGcpAtEnvClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvClientEmail)
	}
	if v := os.Getenv(testAccIntegrationGcpAtEnvResourceID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvResourceID)
	}
	if v := os.Getenv(testAccIntegrationGcpAtEnvSubscription); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpAtEnvSubscription)
	}
}

func testAccIntegrationGcpAtConfig(enabled bool) string {
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
    subscription = "%s"
}
`,
		testAccIntegrationGcpAtResourceType,
		testAccIntegrationGcpAtResourceName,
		enabled,
		os.Getenv(testAccIntegrationGcpAtEnvClientID),
		os.Getenv(testAccIntegrationGcpAtEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGcpAtEnvClientEmail),
		os.Getenv(testAccIntegrationGcpAtEnvPrivateKey),
		os.Getenv(testAccIntegrationGcpAtEnvResourceID),
		resourceLevelAttrOrEmpty(
			os.Getenv(testAccIntegrationGcpAtEnvResourceLevel),
		),
		os.Getenv(testAccIntegrationGcpAtEnvSubscription),
	)
}
