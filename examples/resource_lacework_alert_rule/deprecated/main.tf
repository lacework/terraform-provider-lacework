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
  alert_channels   = var.channels
  severities       = var.severities
  event_categories = var.event_categories
  alert_categories = var.alert_categories
  resource_groups  = [lacework_resource_group.example.id]
}

resource "lacework_resource_group" "example" {
  name     = var.resource_group_name
  type = "AWS"
  group {
    operator = "OR"
    filter {
      filter_name = "filter1"
      field     = "Account"
      operation = "EQUALS"
      value     = ["*"]
    }
  }
}

variable "resource_group_name" {
  type    = string
  default = "Users for Alert Rules Testing example"
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
  default = ["TECHALLY_013F08F1B3FA97E7D54463DECAEEACF9AEA3AEACF863F76"]
}

variable "severities" {
  type    = list(string)
  default = ["High", "Medium"]
}

variable "event_categories" {
  type    = list(string)
  default = ["Compliance", "Platform", "User", "Cloud"]
}

variable "alert_categories" {
  type    = list(string)
  default = ["Policy"]
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

output "alert_categories" {
  value = lacework_alert_rule.example.alert_categories
}

output "resource_group_id" {
  value = lacework_resource_group.example.id
}
