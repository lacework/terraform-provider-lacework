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

  organization {
    admin_accounts = var.admin_accounts
    user_accounts  = var.user_accounts
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
  default = []
}

variable "user_accounts" {
  type    = list(string)
  default = ["MY-ACCOUNT"]
}
