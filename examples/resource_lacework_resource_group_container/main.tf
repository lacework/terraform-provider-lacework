terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group_container" "example" {
  name            = var.resource_group_name
  description     = var.description
  container_tags  = ["myTag"]
  container_labels {
    key = "*"
    value = "*"
  }
}
