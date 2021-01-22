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
	testAccAlertChannelGcpPubSubResourceType = "lacework_alert_channel_gcp_pub_sub"
	testAccAlertChannelGcpPubSubResourceName = "example"

	// Environment variables for testing Gcp Pub Sub Alert Channel Integrations
	testAccAlertChannelGcpPubSubProjectID    = "GCP_PUB_SUB_PROJECT_ID"
	testAccAlertChannelGcpPubSubTopicID      = "GCP_PUB_SUB_TOPIC_ID"
	testAccAlertChannelGcpPubSubClientID     = "GCP_PUB_SUB_CLIENT_ID"
	testAccAlertChannelGcpPubSubClientEmail  = "GCP_PUB_SUB_CLEINT_EMAIL"
	testAccAlertChannelGcpPubSubPrivateKey   = "GCP_PUB_SUB_PRIVATE_KEY"
	testAccAlertChannelGcpPubSubPrivateKeyID = "GCP_PUB_SUB_PRIVATE_KEY_ID"
)

func TestAccAlertChannelGcpPubSub(t *testing.T) {
	resourceTypeAndName := fmt.Sprintf("%s.%s",
		testAccAlertChannelGcpPubSubResourceType,
		testAccAlertChannelGcpPubSubResourceName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlertChannelGcpPubSubEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlertChannelGcpPubSubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannelGcpPubSubConfig(
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelGcpPubSubExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
				),
			},
			{
				Config: testAccAlertChannelGcpPubSubConfig(
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertChannelGcpPubSubExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlertChannelGcpPubSubDestroy(s *terraform.State) error {
	lacework := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccAlertChannelGcpPubSubResourceType {
			continue
		}

		response, err := lacework.Integrations.GetGcpPubSubAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return fmt.Errorf(
					"the %s integration (%s) still exists",
					api.GcpPubSubChannelIntegration, rs.Primary.ID,
				)
			}
		}
	}

	return nil
}

func testAccCheckAlertChannelGcpPubSubExists(resourceTypeAndName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		lacework := testAccProvider.Meta().(*api.Client)

		rs, ok := s.RootModule().Resources[resourceTypeAndName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resourceTypeAndName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource (%s) ID not set", resourceTypeAndName)
		}

		response, err := lacework.Integrations.GetGcpPubSubAlertChannel(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(response.Data) < 1 {
			return fmt.Errorf("the %s integration (%s) doesn't exist",
				api.GcpPubSubChannelIntegration, rs.Primary.ID)
		}

		for _, integration := range response.Data {
			if integration.IntgGuid == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("the %s integration (%s) doesn't exist",
			api.GcpPubSubChannelIntegration, rs.Primary.ID)
	}
}

func testAccAlertChannelGcpPubSubEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(testAccAlertChannelGcpPubSubProjectID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubProjectID)
	}
	if v := os.Getenv(testAccAlertChannelGcpPubSubTopicID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubTopicID)
	}
	if v := os.Getenv(testAccAlertChannelGcpPubSubClientID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubClientID)
	}
	if v := os.Getenv(testAccAlertChannelGcpPubSubClientEmail); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubClientEmail)
	}
	if v := os.Getenv(testAccAlertChannelGcpPubSubPrivateKey); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubPrivateKey)
	}
	if v := os.Getenv(testAccAlertChannelGcpPubSubPrivateKeyID); v == "" {
		t.Fatalf("%s must be set for acceptance tests", testAccAlertChannelGcpPubSubPrivateKeyID)
	}
}

func testAccAlertChannelGcpPubSubConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "integration test"
  enabled = %t
  project_id = "%s"
  topic_id = "%s"
  credentials {
    client_id = "%s"
    client_email = "%s"
    private_key = "%s"
    private_key_id = "%s"
  }
}
		`,
		testAccAlertChannelGcpPubSubResourceType,
		testAccAlertChannelGcpPubSubResourceName,
		enabled,
		os.Getenv(testAccAlertChannelGcpPubSubProjectID),
		os.Getenv(testAccAlertChannelGcpPubSubTopicID),
		os.Getenv(testAccAlertChannelGcpPubSubClientID),
		os.Getenv(testAccAlertChannelGcpPubSubClientEmail),
		os.Getenv(testAccAlertChannelGcpPubSubPrivateKey),
		os.Getenv(testAccAlertChannelGcpPubSubPrivateKeyID),
	)
}
