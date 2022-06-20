terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_policy_exception" "example" {
  policy_id   = var.policyID
  description = var.description
  constraint {
    field_key   = "accountIds"
    field_values = ["*"]
  }
}

variable "policyID" {
  type    = string
  default = "lacework-global-46"
}

variable "description" {
  type    = string
  default = "Policy Exception Created via Terraform"
}

output "lacework_policy_exception" {
  value = lacework_policy_exception.example.description
}
