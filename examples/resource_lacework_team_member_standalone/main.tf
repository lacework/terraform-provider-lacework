terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_team_member" "example" {
  email      = "vatasha.white+2@lacework.net"
  first_name = "Vatasha"
  last_name  = "White"
  company    = "Pokemon International Company"
}
