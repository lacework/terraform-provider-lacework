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
  email      = var.email
  first_name = var.first_name
  last_name  = var.last_name
  company    = "Pokemon International Company"
  enabled    = false

  organization {
    administrator = true
  }
}

variable "email" {
  type    = string
  default = "vatasha.white+1@lacework.net"
}

variable "first_name" {
  type    = string
  default = "Vatasha"
}

variable "last_name" {
  type    = string
  default = "White"
}

variable "admin_accounts" {
  type    = list(string)
  default = ["CUSTOMERDEMO"]
}

variable "user_accounts" {
  type    = list(string)
  default = []
}
