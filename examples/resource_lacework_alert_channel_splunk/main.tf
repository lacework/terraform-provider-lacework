provider "lacework" {}

resource "lacework_alert_channel_splunk" "example" {
  name      = "Splunk Channel Alert Example"
  hec_token = "BA696D5E-CA2F-4347-97CB-3C89F834816F"
  host = "host"
  port = 80
  event_data {
    index = "index"
    source = "source"
  }
}
