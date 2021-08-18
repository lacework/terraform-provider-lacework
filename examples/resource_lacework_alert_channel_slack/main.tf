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
  default = "Slack Alert Channel Example"
}

resource "lacework_alert_channel_slack" "example" {
  name      = var.channel_name
  slack_url = "https://hooks.slack.com/services/ABCD/12345/abcd1234"

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}
