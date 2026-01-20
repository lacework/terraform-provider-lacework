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

resource "lacework_integration_aws_org_agentless_scanning" "example" {
  name                      = var.name
  query_text                = var.query_text
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  scan_multi_volume         = false
  scan_stopped_instances    = true
  scan_short_lived_instances = false
  account_id                = var.account_id
  bucket_arn                = var.bucket_arn
  scanning_account          = var.scanning_account
  management_account        = var.management_account
  monitored_accounts        = var.monitored_accounts

  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }
  
  dynamic "org_account_mappings" {
    for_each = var.org_account_mappings
    content {
      default_lacework_account = org_account_mappings.value["default_lacework_account"]

      dynamic "mapping" {
        for_each = org_account_mappings.value["mapping"]
        content {
          lacework_account = mapping.value["lacework_account"]
          aws_accounts     = mapping.value["aws_accounts"]
        }
      }
    }
  }
}

variable "account_id" {
  type    = string
  default = ""
}

variable "bucket_arn" {
  type    = string
  default = ""
}

variable "role_arn" {
  type    = string
  default = ""
}

variable "external_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = "AWS Organizations Agentless Scanning Example"
}

variable "query_text" {
  type    = string
  default = ""
}

variable "scanning_account" {
  type    = string
  default = ""
}

variable "management_account" {
  type    = string
  default = ""
}

variable "monitored_accounts" {
  type    = list(string)
  default = []
}

variable "org_account_mappings" {
  type = list(object({
    default_lacework_account = string
    mapping = list(object({
      lacework_account = string
      aws_accounts     = list(string)
    }))
  }))
  default     = []
  description = "Mapping of AWS accounts to Lacework accounts within a Lacework organization"
}


output "name" {
  value = lacework_integration_aws_org_agentless_scanning.example.name
}

output "org_account_mappings" {
  value = lacework_integration_aws_org_agentless_scanning.example.org_account_mappings
}

