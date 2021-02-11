provider "lacework" {}

resource "lacework_alert_channel_microsoft_teams" "example" {
  name      = "Microsoft Teams Channel Alert Example"
  webhook_url = "https://outlook.office.com/webhook/api-token"
}
