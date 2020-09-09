provider "lacework" {}

resource "lacework_integration_gcp_cfg" "example" {
  name = "Example-GCP-Integration"
  credentials {
    client_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    private_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_email   = "email@some-project-name.iam.gserviceaccount.com"
    private_key    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  resource_level = "PROJECT"
  resource_id    = "example-project_id"
}
