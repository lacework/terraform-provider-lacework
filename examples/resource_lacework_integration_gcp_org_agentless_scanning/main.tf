terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {
    organization = true
}

variable "name" {
  type    = string
  default = "GCP Agentless Scanning org_example"
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
  default = "org-example-project-id"
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

variable "scan_frequency" {
  type = number 
  default = 24
}

variable "org_account_mappings" {
  type = list(object({
    default_lacework_account = string
    mapping = list(object({
      lacework_account = string
      gcp_projects     = list(string)
    }))
  }))
  default     = []
  description = "Mapping of GCP projects to Lacework accounts within a Lacework organization"
}

resource "lacework_integration_gcp_agentless_scanning" "org_example" {
  name = var.name
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key_id = var.private_key_id
    private_key    = var.private_key
    token_uri 	   = var.token_uri
  }
  resource_level = "ORGANIZATION"
  resource_id    = "techally-test"
  bucket_name = var.bucket_name
  scanning_project_id = "gcp-lw-scanner"
  scan_frequency            = var.scan_frequency
  scan_containers           = true
  scan_host_vulnerabilities = true
  scan_multi_volume         = false
  scan_stopped_instances    = true
  query_text = var.query_text
  filter_list = var.filter_list

  dynamic "org_account_mappings" {
    for_each = var.org_account_mappings
    content {
      default_lacework_account = org_account_mappings.value["default_lacework_account"]

      dynamic "mapping" {
        for_each = org_account_mappings.value["mapping"]
        content {
          lacework_account = mapping.value["lacework_account"]
          gcp_projects     = mapping.value["gcp_projects"]
        }
      }
    }
  }
}

output "name" {
  value = lacework_integration_gcp_agentless_scanning.org_example.name
}

output "client_id" {
  value = lacework_integration_gcp_agentless_scanning.org_example.credentials[0].client_id
}

output "client_email" {
  value = lacework_integration_gcp_agentless_scanning.org_example.credentials[0].client_email
}

output "bucket_name" {
  value = lacework_integration_gcp_agentless_scanning.org_example.bucket_name
}

output "scanning_project_id" {
  value = lacework_integration_gcp_agentless_scanning.org_example.scanning_project_id
}

output "scan_frequency" {
  value = lacework_integration_gcp_agentless_scanning.org_example.scan_frequency
}

output "server_token" {
  value = lacework_integration_gcp_agentless_scanning.org_example.server_token
}

output "org_account_mappings" {
  value = lacework_integration_gcp_agentless_scanning.org_example.org_account_mappings
}