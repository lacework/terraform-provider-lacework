provider "lacework" {}

resource "lacework_integration_docker_v2" "example" {
  name            = "My Docker V2 Registry"
  registry_domain = "my-dockerv2.jfrog.io"
  username        = "my-user"
  password        = "a-secret-password"
  ssl             = true
  limit_by_tags   = ["dev*", "*test"]

  limit_by_labels = {
    key1 = "label1"
    key2 = "label2"
  }
}
