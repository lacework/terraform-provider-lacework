terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_splunk" "example" {
  name      = var.channel_name
  channel   = var.channel
  hec_token = var.hec_token
  host      = var.host
  port      = var.port
  ssl      = var.ssl
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

variable "channel" {
type    = string
default = "Splunk Channel"
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

variable "ssl" {
type    = bool
default = true
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

output "channel" {
  value = lacework_alert_channel_splunk.example.channel
}

output "host" {
  value = lacework_alert_channel_splunk.example.host
}

output "port" {
  value = lacework_alert_channel_splunk.example.port
}

output "ssl" {
  value = lacework_alert_channel_splunk.example.ssl
}