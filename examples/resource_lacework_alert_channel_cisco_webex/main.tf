provider "lacework" {}

resource "lacework_alert_channel_cisco_webex" "example" {
  name       = "My Cisco Webex Channel Alert Example"
  webhook_url = "https://webexapis.com/v1/webhooks/incoming/api-token"
}
