terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_policy" "example" {
  title            = var.title
  query_id         = "LW_Global_AWS_CTA_CloudTrailChange"
  severity         = var.severity
  type             = var.type
  description      = var.description
  remediation      = var.remediation
  evaluation       = var.evaluation
  enabled          = true
  policy_id_suffix = var.policy_id_suffix
  tags = var.tags

  alerting {
    enabled = false
    profile = "LW_CloudTrail_Alerts"
  }
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

variable "type" {
  type    = string
  default = "Violation"
}

variable "evaluation" {
  type    = string
  default = "Hourly"
}

variable "policy_id_suffix" {
  default = ""
}

variable "tags" {
  type = list(string)
  default = ["domain:AWS", "custom"]
}

output "title" {
  value = lacework_policy.example.title
}

output "severity" {
  value = lacework_policy.example.severity
}

output "remediation" {
  value = lacework_policy.example.remediation
}

output "evaluation" {
  value = lacework_policy.example.evaluation
}

output "type" {
  value = lacework_policy.example.type
}

output "description" {
  value = lacework_policy.example.description
}

output "policy_id_suffix" {
  value = lacework_policy.example.policy_id_suffix
}

output "tags" {
  value = lacework_policy.example.tags
}
