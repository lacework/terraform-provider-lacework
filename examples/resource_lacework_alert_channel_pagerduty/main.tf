provider "lacework" {}

resource "lacework_alert_channel_pagerduty" "example" {
  name            = var.channel_name
  integration_key = var.integration_key
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

variable "channel_name" {
  type    = string
  default = "PagerDuty Alert Channel"
}

variable "integration_key" {
  type    = string
  default = "1234abc8901abc567abc123abc78e012"
}

output "channel_name" {
  value = lacework_alert_channel_pagerduty.example.name
}

output "integration_key" {
  value = lacework_alert_channel_pagerduty.example.integration_key
}