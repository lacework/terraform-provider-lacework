provider "lacework" {}

resource "lacework_alert_channel_qradar" "example" {
  name      = "QRadar Channel Alert Example"
  host_url = "https://qradar-lacework.com"
  host_port = 80
  communication_type = "HTTPS"
}
