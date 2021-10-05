terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

variable "channel_name" {
  type    = string
  default = "Datadog Alert Channel Example"
}

variable "api_key" {
  type      = string
  sensitive = true
}

variable "datadog_site" {
  type    = string
  default = "com"
}

variable "datadog_service" {
  type    = string
  default = "Logs Detail"
}

resource "lacework_alert_channel_datadog" "example" {
  name            = var.channel_name
  datadog_site    = var.datadog_site
  datadog_service = var.datadog_service
  api_key         = var.api_key

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}
