terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "name" {
  type    = string
  default = "GCP pub sub audit log integration example"
}

variable "client_id" {
  type    = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "client_email" {
  type    = string
  default = "gcp-al-pub-sub-test@techally-test.iam.gserviceaccount.com"
}

variable "private_key_id" {
  type    = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "private_key" {
  type    = string
  default = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

variable "integration_type" {
  type    = string
  default = "PROJECT"
}

variable "project_id" {
  type    = string
  default = "example-project-id"
}

variable "subscription_id" {
  type    = string
  default = "projects/example-project-id/subscriptions/example-subscription"
}

variable "topic_id" {
  type    = string
  default = "projects/example-project-id/topics/example-topic"
}

resource "lacework_integration_gcp_pub_sub_audit_log" "example" {
  name = var.name
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key_id = var.private_key_id
    private_key    = var.private_key
  }
  integration_type = var.integration_type
  project_id       = var.project_id
  subscription_id  = var.subscription_id
  topic_id         = var.topic_id
  retries          = 10
}

output "name" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.name
}

output "client_id" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.credentials[0].client_id
}

output "client_email" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.credentials[0].client_email
}

output "integration_type" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.integration_type
}

output "project_id" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.project_id
}

output "subscription_id" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.subscription_id
}

output "topic_id" {
  value = lacework_integration_gcp_pub_sub_audit_log.example.topic_id
}

