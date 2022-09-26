terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_gcp_cfg" "example" {
  name = var.integration_name
  credentials {
    client_id      = var.client_id
    client_email   = var.client_email
    private_key_id = var.private_key_id
    private_key    = var.private_key
  }
  resource_level = "PROJECT"
  resource_id    = "techally-test"
  retries        = 10
}


variable "integration_name" {
  type    = string
  default = "Google Cfg Example"
}
variable "client_id" {
  type      = string
  sensitive = true
}
variable "client_email" {
  type      = string
  sensitive = true
}
variable "private_key_id" {
  type      = string
  sensitive = true
}
variable "private_key" {
  type      = string
  sensitive = true
}
variable "non_os_package_support" {
  type    = bool
  default = true
}
