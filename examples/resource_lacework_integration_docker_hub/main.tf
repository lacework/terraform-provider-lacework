provider "lacework" {}

resource "lacework_integration_docker_hub" "example" {
  name           = "My Docker Hub Registry Example"
  username       = "my-user"
  password       = "a-secret-password"
  limit_by_tag   = "dev*"
  limit_by_label = "*label"
  limit_by_repos = "my-repo,other-repo"
  limit_num_imgs = 10
}
