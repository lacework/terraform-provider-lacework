terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_integration_aws_cfg" "example" {
  name = var.name
  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }

  retries = 10
}

variable "name" {
  type = string
  default = "aws cfg created by tf"
}

variable "role_arn" {
  type = string
  default = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
}

variable "external_id" {
  type = string
  default = "12345"
}

output "name" {
  value = lacework_integration_aws_cfg.example.name
}