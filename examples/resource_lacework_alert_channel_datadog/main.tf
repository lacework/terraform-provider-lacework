provider "lacework" {}

variable "channel_name" {
  type = string
  default = "Datadog Channel Alert Example"
}

resource "lacework_alert_channel_datadog" "example" {
  name      = var.channel_name
  datadog_site = "eu"
  datadog_service = "Events Summary"
  api_key = "datadog-key"
}