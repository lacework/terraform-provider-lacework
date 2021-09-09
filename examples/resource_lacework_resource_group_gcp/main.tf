terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group_gcp" "example" {
  name         = var.resource_group_name
  description  = var.description
  organization = var.organization
  projects     = var.projects
}
