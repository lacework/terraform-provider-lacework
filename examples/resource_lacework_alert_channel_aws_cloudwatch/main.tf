provider "lacework" {}

resource "lacework_alert_channel_aws_cloudwatch" "example" {
  name            = "My AWS CloudWatch Alert Channel Example"
  event_bus_arn   = "arn:aws:events:us-west-2:1234567890:event-bus/default"
  group_issues_by = "Resources"
}
