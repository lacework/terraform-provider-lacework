provider "lacework" {}

resource "lacework_alert_channel_splunk" "example" {
  name      = "Splunk Channel Alert Example"
  hec_token = "AA111111-11AA-1AA1-11AA-11111AA1111A"
  host = "host"
  port = 80
  event_data {
    index = "index"
    source = "source"
  }
}