terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_team_member" "example" {
  email         = var.email
  first_name    = var.first_name
  last_name     = var.last_name
  company       = "Marvel Comics"
  administrator = var.administrator
}

variable "email" {
  type    = string
  default = "vatasha.white+2@lacework.net"
}

variable "first_name" {
  type    = string
  default = "Shuri"
}

variable "last_name" {
  type    = string
  default = "White"
}

variable "administrator" {
  type    = bool
  default = false
}
