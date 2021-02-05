provider "lacework" {}

resource "lacework_alert_channel_datadog" "example" {
  name      = "Datadog Channel Alert Example"
  datadog_site = "eu"
  datadog_service = "Events Summary"
  api_key = "datadog-key"
}
