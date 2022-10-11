terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_integration_aws_org_agentless_scanning" "example" {
  name                      = var.name
  query_text                = var.query_text
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  account_id                = var.account_id
  bucket_arn                = var.bucket_arn
  scanning_account          = var.scanning_account
  management_account        = var.management_account
  monitored_accounts        = var.monitored_accounts

  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
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

output "name" {
  value = lacework_integration_aws_org_agentless_scanning.example.name
}
