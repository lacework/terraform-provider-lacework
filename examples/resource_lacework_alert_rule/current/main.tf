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
  alert_subcategories = var.alert_subcategories
  alert_categories = var.alert_categories
  alert_sources    = var.alert_sources
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

variable "alert_subcategories" {
  type    = list(string)
  default = ["Compliance", "Platform", "User", "Cloud Activity"]
}

variable "alert_categories" {
  type    = list(string)
  default = ["Policy"]
}

variable "alert_sources" {
  type    = list(string)
  default = ["AWS"]
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

output "alert_subcategories" {
  value = lacework_alert_rule.example.alert_subcategories
}

output "alert_categories" {
  value = lacework_alert_rule.example.alert_categories
}

output "alert_sources" {
  value = lacework_alert_rule.example.alert_sources
}