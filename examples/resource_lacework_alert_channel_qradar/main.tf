terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_alert_channel_qradar" "example" {
  name               = var.channel_name
  host_url           = var.host_url
  host_port          = var.host_port
  communication_type = var.communication_type
  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = false
}

variable "channel_name" {
  type    = string
  default = "IbmQRadar Alert Channel"
}

variable "host_url" {
  type    = string
  default = "https://qradar-lacework.com"
}

variable "host_port" {
  type    = number
  default = 4000
}

variable "communication_type" {
  type    = string
  default = "HTTPS"
}

output "channel_name" {
  value = lacework_alert_channel_qradar.example.name
}

output "host_url" {
  value = lacework_alert_channel_qradar.example.host_url
}

output "host_port" {
  value = lacework_alert_channel_qradar.example.host_port
}

output "communication_type" {
  value = lacework_alert_channel_qradar.example.communication_type
}