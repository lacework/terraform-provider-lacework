terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_data_export_rule" "example" {
  name             = var.name
  profile_versions =["",""]
  type             = ""
  integration_ids  = [""]
}

variable "name" {
  type    = string
  default = ""
}

output "name" {
  value = lacework_data_export_rule.example.name
}