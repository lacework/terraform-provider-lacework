terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
    organization = true
  }
}

resource "lacework_integration_aws_org_agentless_scanning" "example" {
  name                      = var.name
  query_text                = var.query_text
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  account_id                = "1234556"
  bucket_arn                = var.bucket_arn
  scanning_account          = var.scanning_account
  management_account        = var.management_account
  monitored_accounts        = var.monitored_accounts

  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }

  org_account_mappings {
    default_lacework_account = "lw_account_1"

    mapping {
      lacework_account = "lw_account_2"
      aws_accounts     = ["234556677", "774564564"]
    }

    mapping {
      lacework_account = "lw_account_3"
      aws_accounts     = ["553453453", "934534535"]
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
