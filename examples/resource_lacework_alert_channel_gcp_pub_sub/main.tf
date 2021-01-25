provider "lacework" {}

resource "lacework_alert_channel_gcp_pub_sub" "example" {
  name       = "Gcp Pub Sub example"
  project_id = "my-sample-project-191923"
  topic_id   = "mytopic"
  issue_grouping   = "Events"
  credentials {
    client_id = "client_id"
    client_email = "foo@example.iam.gserviceaccount.com"
    private_key = "priv_key"
    private_key_id = "p_key_id"
  }
}
