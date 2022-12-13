terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "integration_name" {
  type    = string
  default = "GCP Agentless Scanning Example"
}

variable "client_id" {
  type    = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "client_email" {
  type    = string
  default = "email@some-project-name.iam.gserviceaccount.com"
}

variable "private_key_id" {
  type    = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "private_key" {
  type    = string
  default = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

variable "token_uri" {
  type    = string
  default = "https://oauth2.googleapis.com/token"
}

variable "integration_type" {
  type    = string
  default = "PROJECT"
}

variable "project_id" {
  type    = string
  default = "example-project-id"
}

variable "bucket_name" {
  type    = string
  default = "storage bucket id"
}

variable "scanning_project_id" {
  type = string
  default = "scanning-project-id"
}

variable "query_text" {
  type    = string
  default = ""
}

variable "filter_list" {
  type    = list(string)
  default = ["proj1", "proj2"]
}

resource "lacework_integration_gcp_agentless_scanning" "example" {
  name = var.integration_name
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key_id = var.private_key_id
    private_key    = var.private_key
    token_uri 	   = var.token_uri
  }
  resource_level = "PROJECT"
  resource_id    = "techally-test"
  bucket_name = var.bucket_name
  scanning_project_id = "gcp-lw-scanner"
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  query_text = var.query_text
  filter_list = var.filter_list
}

output "name" {
  value = lacework_integration_gcp_agentless_scanning.example.name
}

output "client_id" {
  value = lacework_integration_gcp_agentless_scanning.example.credentials[0].client_id
}

output "client_email" {
  value = lacework_integration_gcp_agentless_scanning.example.credentials[0].client_email
}

output "bucket_name" {
  value = lacework_integration_gcp_agentless_scanning.example.bucket_name
}

output "scanning_project_id" {
  value = lacework_integration_gcp_agentless_scanning.example.scanning_project_id
}

output "scan_frequency" {
  value = lacework_integration_gcp_agentless_scanning.example.scan_frequency
}

output "server_token" {
  value = lacework_integration_gcp_agentless_scanning.example.server_token
}
