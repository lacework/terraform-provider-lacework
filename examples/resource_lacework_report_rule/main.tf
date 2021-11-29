terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_report_rule" "example" {
  name             = var.name
  description      = var.description
  enabled         = true
  severities      = var.severities
  resource_groups = var.resource_groups
  email_alert_channels = var.channels

  aws_compliance_reports {
    pci = var.aws_pci
    cis_s3 = true
  }

  gcp_compliance_reports {
    pci = var.gcp_pci
    cis = true
  }

  daily_compliance_reports {
    aws_cloudtrail = var.daily_cloudtrail
  }

  weekly_snapshot = var.snapshot
}

variable "name" {
  type    = string
  default = "Report Rule 1"
}

variable "description" {
  type    = string
  default = "Report Rule created by Terraform"
}

variable "severities" {
  type    = list(string)
  default = ["High", "Medium"]
}

variable "resource_groups" {
  type    = list(string)
  default = ["TECHALLY_8416B4ADCED28565254842AA5906B729174653E1725F107"]
}

variable "channels" {
  type    = list(string)
  default = ["TECHALLY_2F0C086E17AB64BEC84F4A5FF8A3F068CF2CE15847BCBCA"]
}

variable "aws_pci" {
  type    = bool
  default = true
}

variable "gcp_pci" {
  type    = bool
  default = true
}

variable "daily_cloudtrail" {
  type    = bool
  default = true
}

variable "snapshot" {
  type    = bool
  default = true
}

output "name" {
  value = lacework_report_rule.example.name
}

output "description" {
  value = lacework_report_rule.example.description
}

output "severities" {
  value = lacework_report_rule.example.severities
}

output "resource_groups" {
  value = lacework_report_rule.example.resource_groups
}

output "channels" {
  value = lacework_report_rule.example.email_alert_channels
}

output "aws_pci" {
  value = lacework_report_rule.example.aws_compliance_reports.0.pci
}