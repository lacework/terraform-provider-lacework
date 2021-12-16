provider "lacework" {}

terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_jira_cloud" "example" {
  name                 = var.channel_name
  jira_url             = var.jira_url
  issue_type           = var.issue_type
  project_key          = var.project_key
  username             = var.username
  api_token            = var.api_token
  group_issues_by      = var.group_issues_by
  custom_template_file = var.custom_template_file
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = var.test_integration
}

variable "channel_name" {
  type    = string
  default = "My Jira Cloud Alert Channel Example"
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

variable "api_token" {
  type      = string
  sensitive = true
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
  value = lacework_alert_channel_jira_cloud.example.name
}

output "jira_url" {
  value = lacework_alert_channel_jira_cloud.example.jira_url
}

output "issue_type" {
  value = lacework_alert_channel_jira_cloud.example.issue_type
}

output "project_key" {
  value = lacework_alert_channel_jira_cloud.example.project_key
}

output "group_issues_by" {
  value = lacework_alert_channel_jira_cloud.example.group_issues_by
}

output "custom_template_file" {
  value = lacework_alert_channel_jira_cloud.example.custom_template_file
}

output "username" {
  value = lacework_alert_channel_jira_cloud.example.username
}

output "api_token" {
  value     = lacework_alert_channel_jira_cloud.example.api_token
  sensitive = true
}