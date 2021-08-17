terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_ghcr" "example" {
  name                   = var.integration_name
  registry_notifications = false
  username = var.username
  password = var.password
  ssl      = var.ssl

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["repo/my-image", "repo/other-image"]

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
