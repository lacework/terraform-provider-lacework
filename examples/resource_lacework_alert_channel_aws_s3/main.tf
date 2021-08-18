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

resource "lacework_alert_channel_aws_s3" "example" {
  name       = var.channel_name
  bucket_arn = "arn:aws:s3:::bucket_name/key_name"
  credentials {
    external_id = "12345"
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
  }

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}
