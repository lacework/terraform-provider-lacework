terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

variable "event_bus_arn" {
  description = "The ARN for the event bus"
  type        = string
  sensitive   = true
}

variable "name" {
  description = "The name of the alert channel"
  type        = string
}

resource "lacework_alert_channel_aws_cloudwatch" "example" {
  name            = var.name
  event_bus_arn   = var.event_bus_arn
  group_issues_by = "Events"
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}