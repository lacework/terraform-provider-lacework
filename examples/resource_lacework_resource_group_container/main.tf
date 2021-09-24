terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group_container" "example" {
  name           = var.resource_group_name
  description    = var.description
  container_tags = var.ctr_tags
  container_labels {
    key   = var.ctr_key
    value = var.ctr_value
  }
}
