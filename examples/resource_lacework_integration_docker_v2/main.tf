provider "lacework" {}

resource "lacework_integration_docker_v2" "example" {
  name                   = "My Docker V2 Registry"
  non_os_package_support = true
  registry_domain        = "my-dockerv2.jfrog.io"
  username               = "my-user"
  password               = "a-secret-password"
  ssl                    = true
  notifications          = true
  limit_by_tags          = ["dev*", "*test"]

  limit_by_label {
    key   = "key"
    value = "value"
  }

  limit_by_label {
    key   = "key"
    value = "value2"
  }

  limit_by_label {
    key   = "foo"
    value = "bar"
  }
}
