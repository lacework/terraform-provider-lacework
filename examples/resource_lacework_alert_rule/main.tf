terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_rule" "example" {
  name             = var.name
  description      = var.description
  alert_channels         = var.channels
  severities       = var.severities
  event_categories = var.event_categories
  resource_groups  = [lacework_resource_group_aws.example.id]
}

resource "lacework_resource_group_aws" "example" {
  name     = "Users for Aler Rules Testing"
  accounts = ["*"]
}

variable "name" {
  type    = string
  default = "Alert Rule"
}

variable "description" {
  type    = string
  default = "Alert Rule created by Terraform"
}

variable "channels" {
  type    = list(string)
  default = ["TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350"]
}

variable "severities" {
  type    = list(string)
  default = ["High", "Medium"]
}

variable "event_categories" {
  type    = list(string)
  default = ["Compliance"]
}

output "name" {
  value = lacework_alert_rule.example.name
}

output "description" {
  value = lacework_alert_rule.example.description
}

output "channels" {
  value = lacework_alert_rule.example.alert_channels
}

output "severities" {
  value = lacework_alert_rule.example.severities
}

output "event_categories" {
  value = lacework_alert_rule.example.event_categories
}

output "resource_group_id" {
  value = lacework_resource_group_aws.example.id
}
