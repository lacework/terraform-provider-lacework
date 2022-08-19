terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_data_export_rule" "example" {
  name             = var.name
  integration_ids  = var.integration_ids
}

variable "name" {
  type    = string
  default = "Data Export Rule From Terraform"
}

variable "integration_ids" {
  type    = list(string)
  default = ["TECHALLY_E839836BC385C452E68B3CA7EB45BA0E7BDA39CCF65673A"]
}

output "name" {
  value = lacework_data_export_rule.example.name
}

output "profile_versions" {
  value = lacework_data_export_rule.example.profile_versions
}

output "integration_ids" {
  value = lacework_data_export_rule.example.integration_ids
}

output "type" {
  value = lacework_data_export_rule.example.type
}

output "enabled" {
  value = lacework_data_export_rule.example.enabled
}