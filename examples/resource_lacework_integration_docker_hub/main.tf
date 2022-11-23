terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_docker_hub" "example" {
  name                   = var.integration_name
  non_os_package_support = var.non_os_package_support
  username               = var.user
  password               = var.pass
  limit_num_imgs         = 10
  limit_by_tags          = ["dev*", "*test"]
  limit_by_repositories  = ["my-repo", "other-repo"]

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

variable "integration_name" {
  type    = string
  default = "Dockerhub Container Registry Example"
}
variable "user" {
  type      = string
  sensitive = true
}
variable "pass" {
  type      = string
  sensitive = true
}
variable "non_os_package_support" {
  type    = bool
  default = true
}
