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
  default   = "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks"
}

variable "s3_bucket_arn" {
  type      = string
  default   = "arn:aws:s3:::example-bucket-name"
}

variable "external_id" {
  type      = string
  default   = "12345"
}

variable "role_arn" {
  type      = string
  default   = "arn:aws:iam::249446771485:role/lacework-iam-example-role"
}

resource "lacework_integration_aws_eks_audit_log" "example" {
  name          = var.name
  sns_arn       = var.sns_arn
  s3_bucket_arn = var.s3_bucket_arn
  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }
  retries = 10
}

output "name" {
  value = lacework_integration_aws_eks_audit_log.example.name
}

output "sns_arn" {
  value = lacework_integration_aws_eks_audit_log.example.sns_arn
}

output "s3_bucket_arn" {
  value = lacework_integration_aws_eks_audit_log.example.s3_bucket_arn
}

output "role_arn" {
  value = lacework_integration_aws_eks_audit_log.example.credentials[0].role_arn
}

output "external_id" {
  value = lacework_integration_aws_eks_audit_log.example.credentials[0].external_id
}
