provider "lacework" {}

resource "lacework_alert_channel_datadog" "example" {
  name      = "Datadog Channel Alert Example"
  teams_url = "eu"
}
