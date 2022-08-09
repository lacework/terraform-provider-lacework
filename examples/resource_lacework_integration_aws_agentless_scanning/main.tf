terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_integration_aws_agentless_scanning" "example" {
  name                      = var.name
  query_text                = var.query_text
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  account_id = var.account_id
  bucket_arn = var.bucket_arn
  credentials {
    role_arn = var.role_arn
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
  default = "AWS Agentless Scanning Example"
}

variable "query_text" {
  type    = string
  default = ""
}

output "name" {
  value = lacework_integration_aws_agentless_scanning.example.name
}