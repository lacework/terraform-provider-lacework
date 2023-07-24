terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_managed_policies" "example" {
  policy {
    id       = var.id_1
    enabled  = var.enabled_1
    severity = var.severity_1
  }
  policy {
    id       = var.id_2
    enabled  = var.enabled_2
    severity = var.severity_2
  }
}

variable "id_1" {
  type    = string
  default = "lacework-global-1"
}

variable "enabled_1" {
  type    = bool
  default = true
}

variable "severity_1" {
  type    = string
}

variable "id_2" {
  type    = string
  default = "lacework-global-2"
}

variable "enabled_2" {
  type    = bool
  default = true
}

variable "severity_2" {
  type    = string
}

output "policy" {
  value = lacework_managed_policies.example.policy
}
