terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_alert_channel_microsoft_teams" "example" {
  name        = var.channel_name
  webhook_url = var.webhook_url

  // test_integration input is used in this example only for testing
  // purposes, it help us avoid sending a "test" request to the
  // system we are integrating to. In production, this should remain
  // turned on ("true") which is the default setting
  test_integration = var.test_integration
}

// Variables and Outputs used in our integration tests (integration/)
variable "test_integration" {
  type    = bool
  default = false
}

variable "channel_name" {
  type    = string
  default = "Microsoft Teams Alert Channel Example"
}

variable "webhook_url" {
  type    = string
  default = "https://test-tenant.outlook.office.com/webhook/api-token"
}

output "channel_name" {
  value = lacework_alert_channel_microsoft_teams.example.name
}

output "webhook_url" {
  value = lacework_alert_channel_microsoft_teams.example.webhook_url
}
