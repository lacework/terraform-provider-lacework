terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_policy_compliance" "example" {
  title            = var.title
  query_id         = "LW_Global_AWS_Config_S3BucketLoggingNotEnabled"
  severity         = var.severity
  description      = var.description
  remediation      = var.remediation
  enabled          = true
  policy_id_suffix = var.policy_id_suffix
  tags             = var.tags
  alerting_enabled = true
}


variable "title" {
  type    = string
  default = "lql-terraform-policy"
}

variable "description" {
  type    = string
  default = "Policy Created via Terraform"
}

variable "severity" {
  type    = string
  default = "High"
}

variable "remediation" {
  type    = string
  default = "Please Investigate"
}

variable "policy_id_suffix" {
  default = ""
}

variable "tags" {
  type    = list(string)
  default = ["cloud_AWS", "custom"]
}

output "title" {
  value = lacework_policy_compliance.example.title
}

output "severity" {
  value = lacework_policy_compliance.example.severity
}

output "remediation" {
  value = lacework_policy_compliance.example.remediation
}

output "description" {
  value = lacework_policy_compliance.example.description
}

output "policy_id_suffix" {
  value = lacework_policy_compliance.example.policy_id_suffix
}

output "tags" {
  value = lacework_policy_compliance.example.tags
}
