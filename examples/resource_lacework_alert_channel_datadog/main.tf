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
  default = "Datadog Alert Channel Test Example 2"
}

resource "lacework_alert_channel_datadog" "example" {
  name             = var.channel_name
  datadog_site     = "eu"
  datadog_service  = "Events Summary"
  api_key          = "datadog-key"

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}
