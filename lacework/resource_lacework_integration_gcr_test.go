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
	testAccIntegrationGcrResourceType = "lacework_integration_gcr"
	testAccIntegrationGcrResourceName = "example"

	// Environment variables for testing GCR Integrations
	testAccIntegrationGcrEnvRegistryDomain = "GCR_DOMAIN"
	testAccIntegrationGcrEnvClientID       = "GCR_CLIENT_ID"
	testAccIntegrationGcrEnvPrivateKeyID   = "GCR_PRIVATE_KEY_ID"
	testAccIntegrationGcrEnvPrivateKey     = "GCR_PRIVATE_KEY"
	testAccIntegrationGcrEnvClientEmail    = "GCR_CLIENT_EMAIL"
)

func TestAccIntegrationGcr(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationGcrResourceType,
		testAccIntegrationGcrResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationGcrEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationGcrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationGcrConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcrExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationGcrConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationGcrExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationGcrDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationGcrResourceType {
			continue
		}

		response, err := lacework.Integrations.GetSlackAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.SlackChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckIntegrationGcrExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetSlackAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.SlackChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.SlackChannelIntegration, rs.Primary.ID)
	}
}

func testAccIntegrationGcrEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitByTag); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitByTag)
	}
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitByLabel); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitByLabel)
	}
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitByRepos); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitByRepos)
	}
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitNumImg); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitNumImg)
	}
	if v := os.Getenv(testAccIntegrationGcrEnvRegistryDomain); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcrEnvRegistryDomain)
	}
	if v := os.Getenv(testAccIntegrationGcrEnvClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcrEnvClientID)
	}
	if v := os.Getenv(testAccIntegrationGcrEnvPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcrEnvPrivateKeyID)
	}
	if v := os.Getenv(testAccIntegrationGcrEnvPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcrEnvPrivateKey)
	}
	if v := os.Getenv(testAccIntegrationGcrEnvClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationGcrEnvClientEmail)
	}
}

func testAccIntegrationGcrConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name            = "integration test"
    enabled         = %t
    registry_domain = "%s"
    credentials {
        client_id = "%s"
        private_key_id = "%s"
        client_email = "%s"
        private_key = "%s"
    }
    limit_by_tag    = "%s"
    limit_by_label  = "%s"
    limit_by_repos = "%s"
    limit_num_imgs = %s
}
`,
		testAccIntegrationGcrResourceType,
		testAccIntegrationGcrResourceName,
		enabled,
		os.Getenv(testAccIntegrationGcrEnvRegistryDomain),
		os.Getenv(testAccIntegrationGcrEnvClientID),
		os.Getenv(testAccIntegrationGcrEnvPrivateKeyID),
		os.Getenv(testAccIntegrationGcrEnvClientEmail),
		os.Getenv(testAccIntegrationGcrEnvPrivateKey),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByTag),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByLabel),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByRepos),
		os.Getenv(testAccIntegrationContainerRegEnvLimitNumImg),
	)
}
