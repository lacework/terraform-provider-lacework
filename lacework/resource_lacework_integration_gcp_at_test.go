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

	// Environment variables for testing GCP AT
	testAccIntegrationGcpEnvSubscription = "GCP_SUBSCRIPTION"
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
	if v := os.Getenv(testAccIntegrationGcpEnvResourceID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvResourceID)
	}
	if v := os.Getenv(testAccIntegrationGcpEnvSubscription); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcpEnvSubscription)
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
		os.Getenv(testAccIntegrationGcpEnvClientID),
		os.Getenv(testAccIntegrationGcpEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGcpEnvClientEmail),
		os.Getenv(testAccIntegrationGcpEnvPrivateKey),
		os.Getenv(testAccIntegrationGcpEnvResourceID),
		resourceLevelAttrOrEmpty(
			os.Getenv(testAccIntegrationGcpEnvResourceLevel),
		),
		os.Getenv(testAccIntegrationGcpEnvSubscription),
	)
}
