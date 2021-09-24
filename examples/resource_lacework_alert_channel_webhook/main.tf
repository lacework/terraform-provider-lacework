provider "lacework" {}

resource "lacework_alert_channel_webhook" "example" {
  name        = "My Webhook Channel Alert Example"
  webhook_url = "https://hook.com/webhook?api-token=123"
}
