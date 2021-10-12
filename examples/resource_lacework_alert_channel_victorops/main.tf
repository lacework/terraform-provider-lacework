terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_alert_channel_victorops" "example" {
  name        = var.channel_name
  webhook_url = var.webhook_url
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

output "channel_name" {
  value = lacework_alert_channel_victorops.example.name
}

output "webhook_url" {
  value = lacework_alert_channel_victorops.example.webhook_url
}