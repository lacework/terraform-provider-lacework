terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group_machine" "example" {
  name            = var.resource_group_name
  description     = var.description
  machine_tags {
    key   = var.machine_key
    value = var.machine_value
  }
}
