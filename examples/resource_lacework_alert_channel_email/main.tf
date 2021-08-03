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
  default = "Email Alert Channel Example"
}

resource "lacework_alert_channel_email" "example" {
  name       = var.channel_name
  recipients = ["foo@example.com"]
  test_integration = false
}
