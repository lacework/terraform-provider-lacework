provider "lacework" {}

resource "lacework_integration_docker_v2" "example" {
  name            = "My Docker V2 Registry Example"
  registry_domain = "127.0.0.1:1234"
  username        = "my-user"
  password        = "a-secret-password"
  ssl             = true

  # These attributes are deprecated and will be removed in v1.0
  limit_by_tag   = "dev*"
  limit_by_label = "*label"
}
