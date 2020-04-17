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
	testAccIntegrationAzureActivityLogResourceType = "lacework_integration_azure_al"
	testAccIntegrationAzureActivityLogResourceName = "example"

	// Environment variables for testing AZURE_AL_SEQ specifically
	testAccIntegrationAzureEnvQueueUrl = "AZURE_QUEUE_URL"
)

func TestAccIntegrationAzureActivityLog(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationAzureActivityLogResourceType,
		testAccIntegrationAzureActivityLogResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationAzureActivityLogEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationAzureActivityLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationAzureActivityLogConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAzureActivityLogExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationAzureActivityLogConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAzureActivityLogExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAzureActivityLogDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationAzureActivityLogResourceType {
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

func testAccCheckIntegrationAzureActivityLogExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationAzureActivityLogEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationAzureEnvQueueUrl); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureEnvQueueUrl)
	}
	if v := os.Getenv(testAccIntegrationAzureEnvTenantID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureEnvTenantID)
	}
	if v := os.Getenv(testAccIntegrationAzureEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationAzureEnvClientSecret); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationAzureEnvClientSecret)
	}
}

func testAccIntegrationAzureActivityLogConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "integration test"
    enabled = %t
    tenant_id = "%s"
    queue_url = "%s"
    credentials {
        client_id = "%s"
        client_secret = "%s"
    }
}
`,
		testAccIntegrationAzureActivityLogResourceType,
		testAccIntegrationAzureActivityLogResourceName,
		enabled,
		os.Getenv(testAccIntegrationAzureEnvTenantID),
		os.Getenv(testAccIntegrationAzureEnvQueueUrl),
		os.Getenv(testAccIntegrationAzureEnvClientID),
		os.Getenv(testAccIntegrationAzureEnvClientSecret),
	)
}
