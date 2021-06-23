provider "lacework" {}

resource "lacework_integration_docker_hub" "example" {
  name                  = "My Docker Hub Registry Example"
  username              = "my-user"
  password              = "a-secret-password"
  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["my-repo", "other-repo"]

  limit_by_labels = {
    key1 = "label1"
    key2 = "label2"
  }
}
