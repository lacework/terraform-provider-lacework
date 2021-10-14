terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_splunk" "example" {
  name      = var.channel_name
  hec_token = var.hec_token
  host      = var.host
  port      = var.port
  event_data {
    index  = var.index
    source = var._source
  }
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

variable "channel_name" {
type    = string
default = "Splunk Alert Channel"
}

variable "hec_token" {
type    = string
default = "BA696D5E-CA2F-4347-97CB-3C89F834816F"
}

variable "host" {
type    = string
default = "host"
}

variable "port" {
type    = number
default = 80
}

variable "index" {
type    = string
default = "index"
}

variable "_source" {
type    = string
default = "source"
}

output "channel_name" {
  value = lacework_alert_channel_splunk.example.name
}

output "hec_token" {
  value = lacework_alert_channel_splunk.example.hec_token
}