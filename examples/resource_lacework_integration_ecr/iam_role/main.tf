terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_ecr" "iam_role" {
  name            = var.integration_name
  registry_domain = var.registry_domain
  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }
  non_os_package_support = var.non_os_packages
}
