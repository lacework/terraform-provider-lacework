terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_newrelic" "example" {
  name       = var.channel_name
  account_id = var.account_id
  insert_key = var.insert_key
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

variable "channel_name" {
  type    = string
  default = "NewRelic Insights Channel Alert Example"
}

variable "account_id" {
  type    = number
  default = 2338053
}

variable "insert_key" {
  type    = string
  default = "x-xx-xxxxxxxxxxxxxxxxxx"
}

output "channel_name" {
  value = lacework_alert_channel_newrelic.example.name
}

output "account_id" {
  value = lacework_alert_channel_newrelic.example.account_id
}

output "insert_key" {
  value = lacework_alert_channel_newrelic.example.insert_key
}