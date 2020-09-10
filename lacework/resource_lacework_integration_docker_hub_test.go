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
	testAccIntegrationDockerHubResourceType = "lacework_integration_docker_hub"
	testAccIntegrationDockerHubResourceName = "example"

	// Environment variables for testing Docker Hub Container Registry Integrations
	testAccIntegrationContainerRegEnvLimitByTag   = "LIMIT_BY_TAG"
	testAccIntegrationContainerRegEnvLimitByLabel = "LIMIT_BY_LABEL"
	testAccIntegrationContainerRegEnvLimitByRepos = "LIMIT_BY_REP"
	testAccIntegrationContainerRegEnvLimitNumImg  = "LIMIT_NUM_IMG"
	testAccIntegrationDockerHubEnvUsername        = "DOCKER_USERNAME"
	testAccIntegrationDockerHubEnvPassword        = "DOCKER_PASSWORD"
)

func TestAccIntegrationDockerHub(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccIntegrationDockerHubResourceType,
		testAccIntegrationDockerHubResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccIntegrationDockerHubEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIntegrationDockerHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationDockerHubConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationDockerHubExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccIntegrationDockerHubConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationDockerHubExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIntegrationDockerHubDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccIntegrationDockerHubResourceType {
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

func testAccCheckIntegrationDockerHubExists(resourceTypeAndName string) resource.TestCheckFunc {
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

func testAccIntegrationDockerHubEnvVarsPreCheck(t *testing.T) {
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
	if v := os.Getenv(testAccIntegrationDockerHubEnvUsername); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationDockerHubEnvUsername)
	}
	if v := os.Getenv(testAccIntegrationDockerHubEnvPassword); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccIntegrationDockerHubEnvPassword)
	}
}

func testAccIntegrationDockerHubConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name           = "integration test"
    enabled        = %t
    username       = "%s"
    password       = "%s"
    limit_by_tag   = "%s"
    limit_by_label = "%s"
    limit_by_repos = "%s"
    limit_num_imgs = %s
}
`,
		testAccIntegrationDockerHubResourceType,
		testAccIntegrationDockerHubResourceName,
		enabled,
		os.Getenv(testAccIntegrationDockerHubEnvUsername),
		os.Getenv(testAccIntegrationDockerHubEnvPassword),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByTag),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByLabel),
		os.Getenv(testAccIntegrationContainerRegEnvLimitByRepos),
		os.Getenv(testAccIntegrationContainerRegEnvLimitNumImg),
	)
}
