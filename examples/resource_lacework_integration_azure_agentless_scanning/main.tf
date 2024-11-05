terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "integration_name" {
  type    = string
  default = "Agentless_Scanning_Example"
}

variable "client_id" {
  type    = string
}

variable "client_secret" {
  type    = string
  sensitive = true
}

variable "integration_level" {
  type    = string
  default = "TENANT"
}

variable "blob_container_name" {
  type    = string
  default = "terraform-provider-test"
}

variable "scanning_subscription_id" {
  type = string
  default = "0252a545-04d4-4262-a82c-ceef83344237"
}

variable "scanning_resource_group_name" {
  type = string
  default = "agentless-canary-scanned-group"
}

variable "storage_account_url" {
  type = string
  default = "https://asidekicktest3954.blob.core.windows.net/"
}

variable "tenant_id" {
  type = string
  default = "a329d4bf-4557-4ccf-b132-84e7025ea22d"
}

variable "query_text" {
  type    = string
  default = ""
}

variable "subscriptions_list" {
  type    = list(string)
  default = []
}

resource "lacework_integration_azure_agentless_scanning" "example" {
  name = var.integration_name
  credentials {
    client_id      = var.client_id
    client_secret = var.client_secret
  }
  integration_level = var.integration_level
  blob_container_name = var.blob_container_name
  scanning_subscription_id = var.scanning_subscription_id
  tenant_id = var.tenant_id
  scanning_resource_group_name = var.scanning_resource_group_name
  storage_account_url = var.storage_account_url
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  scan_multi_volume         = false
  scan_stopped_instances    = true
  query_text = var.query_text
  subscriptions_list = var.subscriptions_list
}

output "name" {
  value = lacework_integration_azure_agentless_scanning.example.name
}

output "client_id" {
  value = lacework_integration_azure_agentless_scanning.example.credentials[0].client_id
}

output "blob_container_name" {
  value = lacework_integration_azure_agentless_scanning.example.blob_container_name
}

output "scanning_subscription_id" {
  value = lacework_integration_azure_agentless_scanning.example.scanning_subscription_id
}

output "tenant_id" {
  value = lacework_integration_azure_agentless_scanning.example.tenant_id
}

output "scan_frequency" {
  value = lacework_integration_azure_agentless_scanning.example.scan_frequency
}

output "server_token" {
  value = lacework_integration_azure_agentless_scanning.example.server_token
}
