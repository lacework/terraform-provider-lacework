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
	testAccIntegrationDockerV2ResourceType = "lacework_integration_docker_v2"
	testAccIntegrationDockerV2ResourceName = "example"

	// Environment variables for testing Docker V2 Container Registry Integrations
	testAccIntegrationDockerV2EnvRegistryDomain = "DOCKER_V2_DOMAIN"
	testAccIntegrationDockerV2EnvUsername       = "DOCKER_V2_USERNAME"
	testAccIntegrationDockerV2EnvPassword       = "DOCKER_V2_PASSWORD"
	testAccIntegrationDockerV2EnvSSL            = "DOCKER_V2_SSL"
)

func TestAccIntegrationDockerV2(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationDockerV2ResourceType,
		testAccIntegrationDockerV2ResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationDockerV2EnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationDockerV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationDockerV2Config(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationDockerV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationDockerV2Config(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationDockerV2Exists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationDockerV2Destroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationDockerV2ResourceType {
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

func testAccCheckIntegrationDockerV2Exists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationDockerV2EnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitByTag); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitByTag)
	}
	if v := os.Getenv(testAccIntegrationContainerRegEnvLimitByLabel); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationContainerRegEnvLimitByLabel)
	}
	if v := os.Getenv(testAccIntegrationDockerV2EnvRegistryDomain); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationDockerV2EnvRegistryDomain)
	}
	if v := os.Getenv(testAccIntegrationDockerV2EnvUsername); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationDockerV2EnvUsername)
	}
	if v := os.Getenv(testAccIntegrationDockerV2EnvPassword); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationDockerV2EnvPassword)
	}
}

func testAccIntegrationDockerV2Config(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name            = "integration test"
    enabled         = %t
    registry_domain = "%s"
    username        = "%s"
    password        = "%s"
    ssl             = %s
    limit_by_tag    = "%s"
    limit_by_label  = "%s"
}
`,
		testAccIntegrationDockerV2ResourceType,
		testAccIntegrationDockerV2ResourceName,
		enabled,
		os.Getenv(testAccIntegrationDockerV2EnvRegistryDomain),
		os.Getenv(testAccIntegrationDockerV2EnvUsername),
		os.Getenv(testAccIntegrationDockerV2EnvPassword),
		os.Getenv(testAccIntegrationDockerV2EnvSSL),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByTag),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByLabel),
	)
}
