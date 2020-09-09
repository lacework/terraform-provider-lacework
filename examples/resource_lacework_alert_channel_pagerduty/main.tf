provider "lacework" {}

resource "lacework_alert_channel_pagerduty" "example" {
  name            = "default alerts"
  integration_key = "1234abc8901abc567abc123abc78e012"
}
