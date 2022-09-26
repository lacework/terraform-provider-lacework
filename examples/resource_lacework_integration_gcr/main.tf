terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_gcr" "example" {
  name            = var.integration_name
  registry_domain = "gcr.io"
  non_os_package_support = true
  credentials {
    client_id      = var.client_id
    private_key_id = var.private_key_id
    client_email   = var.client_email
    private_key    = var.private_key
  }

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["my-repo", "other-repo"]

  limit_by_labels {
    key   = "foo"
    value = "bar"
  }
}

variable "integration_name" {
  type    = string
  default = "Google Container Registry Example"
}
variable "client_id" {
  type      = string
  sensitive = true
}
variable "client_email" {
  type      = string
  sensitive = true
}
variable "private_key_id" {
  type      = string
  sensitive = true
}
variable "private_key" {
  type      = string
  sensitive = true
}

