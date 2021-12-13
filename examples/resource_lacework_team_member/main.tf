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
  email         = "vatasha.white+1@lacework.net"
  first_name    = "Vatasha"
  last_name     = "White"
  company       = "Pokemon International Company"
  enabled       = true
  administrator = false
  organization {
    admin_accounts = ["tech-ally", "xyz"]
    user_accounts = ["customerdemo"]

  }
}
