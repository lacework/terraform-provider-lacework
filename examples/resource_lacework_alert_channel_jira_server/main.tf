provider "lacework" {}

terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_jira_server" "example" {
  name                 = var.channel_name
  jira_url             = var.jira_url
  issue_type           = var.issue_type
  project_key          = var.project_key
  username             = var.username
  configuration        = var.configuration
  group_issues_by      = var.group_issues_by
  custom_template_file = var.custom_template_file
  password             = var.password
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

variable "channel_name" {
  type    = string
  default = "My Jira Server Alert Channel Example"
}

variable "jira_url" {
  type = string
}

variable "issue_type" {
  type = string
}

variable "project_key" {
  type = string
}

variable "username" {
  type = string
}

variable "password" {
  type      = string
  sensitive = true
}

variable "configuration" {
  type = string
}

variable "group_issues_by" {
  type = string
}

variable "custom_template_file" {
  type    = string
  default = <<TEMPLATE
{
    "fields": {
        "labels": [
            "myLabel"
        ],
        "priority":
        {
            "id": "1"
        }
    }
}
TEMPLATE
}

variable "test_integration" {
  type    = bool
  default = false
}

output "channel_name" {
  value = lacework_alert_channel_jira_server.example.name
}

output "jira_url" {
  value = lacework_alert_channel_jira_server.example.jira_url
}

output "issue_type" {
  value = lacework_alert_channel_jira_server.example.issue_type
}

output "project_key" {
  value = lacework_alert_channel_jira_server.example.project_key
}

output "configuration" {
  value = lacework_alert_channel_jira_server.example.configuration
}

output "group_issues_by" {
  value = lacework_alert_channel_jira_server.example.group_issues_by
}

output "custom_template_file" {
  value = lacework_alert_channel_jira_server.example.custom_template_file
}

output "username" {
  value = lacework_alert_channel_jira_server.example.username
}

output "password" {
  value     = lacework_alert_channel_jira_server.example.password
  sensitive = true
}