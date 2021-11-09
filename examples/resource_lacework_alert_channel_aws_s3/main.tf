terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

variable "channel_name" {
  type    = string
  default = "AwsS3 Alert Channel Example"
}

variable "bucket_arn" {
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

resource "lacework_alert_channel_aws_s3" "example" {
  name       = var.channel_name
  bucket_arn = var.bucket_arn
  credentials {
    external_id = var.external_id
    role_arn    = var.role_arn
  }

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}
