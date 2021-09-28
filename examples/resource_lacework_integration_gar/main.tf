terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_gar" "example" {
  name            = var.integration_name
  registry_domain = "us-west1-docker.pkg.dev"
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key_id = var.private_key_id
    private_key    = var.private_key
  }
  non_os_package_support = var.non_os_packages

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
