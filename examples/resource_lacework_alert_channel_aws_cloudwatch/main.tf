terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_alert_channel_aws_cloudwatch" "example" {
  name            = "My AWS CloudWatch Alert Channel Example"
  event_bus_arn   = "arn:aws:events:us-west-2:1234567890:event-bus/default"
  group_issues_by = "Events"
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}