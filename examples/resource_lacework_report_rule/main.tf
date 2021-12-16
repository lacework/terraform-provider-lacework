terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_email" "email_alerts" {
  name       = "Used for Report Rules Testing"
  recipients = ["foo@example.com"]

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

resource "lacework_resource_group_aws" "aws_group" {
  name     = var.resource_group_name
  accounts = ["*"]
}

resource "lacework_report_rule" "example" {
  name             = var.name
  description      = var.description
  enabled         = true
  severities      = var.severities
  resource_groups = [lacework_resource_group_aws.aws_group.id]
  email_alert_channels = [lacework_alert_channel_email.email_alerts.id]

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
  default = "Terraform Report Rule"
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

variable "resource_group_name" {
  type    = string
  default = "Used for Report Rules Testing"
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
