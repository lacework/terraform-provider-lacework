terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {
  organization = true
}

resource "lacework_team_member" "example" {
  email         = "vatasha.white+terraformtesting2@lacework.net"
  first_name    = "Vatasha"
  last_name     = "White"
  company       = "Pokemon International Company"
  enabled       = true
  administrator = false
  organization {
    administrator = false
    user = true
    admin_accounts = ["tech-ally", "xyz"]
  }
}
