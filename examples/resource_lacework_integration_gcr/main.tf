provider "lacework" {}

resource "lacework_integration_gcr" "example" {
  name            = "GRC Example"
  registry_domain = "gcr.io"
  credentials {
    client_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    private_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_email   = "email@some-project-name.iam.gserviceaccount.com"
    private_key    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  limit_by_tag   = "dev*"
  limit_by_label = "*label"
  limit_by_repos = "my-repo,other-repo"
  limit_num_imgs = 10
}
