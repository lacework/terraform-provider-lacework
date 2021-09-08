terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group_aws" "example" {
  name        = var.resource_group_name
  description = var.description
  accounts    = var.accounts
}
