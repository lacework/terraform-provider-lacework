terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

variable "name" {
  type    = string
  default = "My GCP Pub Sub Example"
}

variable "project_id" {
  type = string
}

variable "topic_id" {
  type = string
}

variable "issue_grouping" {
  type = string
}

variable "client_id" {
  type = string
}

variable "client_email" {
  type = string
}

variable "private_key" {
  type      = string
  sensitive = true
}

variable "private_key_id" {
  sensitive = true
  type = string
}

resource "lacework_alert_channel_gcp_pub_sub" "example" {
  name           = var.name
  project_id     = var.project_id
  topic_id       = var.topic_id
  issue_grouping = var.issue_grouping
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key    = var.private_key
    private_key_id = var.private_key_id
  }
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

output "name" {
  value = lacework_alert_channel_gcp_pub_sub.example.name
}

output "project_id" {
  value = lacework_alert_channel_gcp_pub_sub.example.project_id
}

output "topic_id" {
  value = lacework_alert_channel_gcp_pub_sub.example.topic_id
}

output "issue_grouping" {
  value = lacework_alert_channel_gcp_pub_sub.example.issue_grouping
}

output "client_id" {
  value = lacework_alert_channel_gcp_pub_sub.example.credentials[0].client_id
}

output "client_email" {
  value = lacework_alert_channel_gcp_pub_sub.example.credentials[0].client_email
}

output "private_key" {
  value     = lacework_alert_channel_gcp_pub_sub.example.credentials[0].private_key
  sensitive = true
}

output "private_key_id" {
  value = lacework_alert_channel_gcp_pub_sub.example.credentials[0].private_key_id
  sensitive = true
}
