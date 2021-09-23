terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {
  organization = true
}

resource "lacework_resource_group_lw_account" "example" {
  name            = var.resource_group_name
  description     = var.description
  accounts        = ["tech-ally"]
}
