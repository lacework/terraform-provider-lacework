provider "lacework" {}

resource "lacework_alert_channel_slack" "example" {
  name      = "My Slack Channel Alert Example"
  slack_url = "https://hooks.slack.com/services/ABCD/12345/abcd1234"
}
