terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_profile" "custom_profile" {
  name    = var.name
  extends = var.extends

  alert {
    name        = "Violation"
    event_name  = "LW Configuration GCP Violation Alert"
    subject     = "{{_OCCURRENCE}} violation detected in project {{PROJECT_ID}}"
    description = var.alert_description
  }
}

variable "name" {
  type    = string
  default = "CUSTOM_PROFILE_TERRAFORM_TEST"
}

variable "extends" {
  type    = string
  default = "LW_CFG_GCP_DEFAULT_PROFILE"
}

variable "alert_description" {
  type    = string
  default = "{{_OCCURRENCE}} violation for GCP Resource {{RESOURCE_TYPE}}:{{RESOURCE_ID}} in project {{PROJECT_ID}} region {{RESOURCE_REGION}}"
}

output "name" {
  value = lacework_alert_profile.custom_profile.name
}

output "extends" {
  value = lacework_alert_profile.custom_profile.extends
}

output "alert_description" {
  value = lacework_alert_profile.custom_profile.alert.*.description[0]
}