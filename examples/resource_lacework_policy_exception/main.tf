terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_policy_exception" "example" {
  policy_id   = var.policy_id
  description = var.description
  constraint {
    field_key    = var.field_key
    field_values = var.field_values
  }
}

variable "policy_id" {
  type    = string
  default = "lacework-global-46"
}
variable "field_key" {
  type    = string
  default = "accountIds"
}
variable "field_values" {
  type    = list(string)
  default = ["*"]
}

variable "description" {
  type    = string
  default = "Policy Exception Created via Terraform"
}

output "description" {
  value = lacework_policy_exception.example.description
}

output "policy_id" {
  value = lacework_policy_exception.example.policy_id
}
