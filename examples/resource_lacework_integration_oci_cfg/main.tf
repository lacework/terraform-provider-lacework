terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_oci_cfg" "example" {
  name = var.name
  credentials {
    fingerprint = var.fingerprint
    private_key = var.private_key
  }
  home_region = var.home_region
  tenant_id = var.tenant_id
  tenant_name = var.tenant_name
  user_ocid = var.user_ocid
  retries = 10
}

variable "name" {
  type = string
  default = "OCI config integration example"
}

variable "fingerprint" {
  type    = string
}

variable "private_key" {
  type    = string
  sensitive = true
}

variable "home_region" {
  type    = string
  default = "us-sanjose-1"
}

variable "tenant_id" {
  type    = string
}

variable "tenant_name" {
  type    = string
}

variable "user_ocid" {
  type    = string
}

output "name" {
  value = lacework_integration_oci_cfg.example.name
}
