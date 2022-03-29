terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "name" {
  type    = string
  default = "AWS EKS audit log integration example"
}

variable "sns_arn" {
  type      = string
  sensitive = true
}

variable "external_id" {
  type      = string
  sensitive = true
}

variable "role_arn" {
  type      = string
}

resource "lacework_integration_aws_eks_audit_log" "example" {
  name      = var.name
  sns_arn = var.sns_arn
  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }
  retries = 10
}
