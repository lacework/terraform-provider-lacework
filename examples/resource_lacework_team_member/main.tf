terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_team_member" "example" {
  email         = "vatasha.white+terraformtest@lacework.net"
  first_name    = "Vatasha"
  last_name     = "White"
  company       = "Pokemon International Company"
  enabled       = false
  administrator = false
}
