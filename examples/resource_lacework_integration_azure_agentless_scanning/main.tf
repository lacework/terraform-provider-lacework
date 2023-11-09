terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "integration_name" {
  type    = string
  default = "Agentless Scanning Example"
}

variable "client_id" {
  type    = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "client_secret" {
  type    = string
  default = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

variable "integration_level" {
  type    = string
  default = "SUBSCRIPTION"
}

variable "blob_container_name" {
  type    = string
  default = "blob container name"
}

# TODO: make this fully qualified 
variable "scanning_subscription_id" {
  type = string
  default = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "scanning_resource_group_id" {
  type = string
  default = "abcd"
}

variable "storage_account_url" {
  type = string
  default = "https://blobabc.blob.core.windows.net"
}

variable "tenant_id" {
  type = string
  default = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

variable "query_text" {
  type    = string
  default = ""
}

variable "subscription_list" {
  type    = list(string)
  default = ["sub1", "sub2"]
}

resource "lacework_integration_azure_agentless_scanning" "example" {
  name = var.integration_name
  credentials {
    client_id      = var.client_id
    client_secret = var.client_secret
  }
  integration_level = "SUBSCRIPTION"
  blob_container_name = var.blob_container_name
  scanning_subscription_id = var.scanning_subscription_id
  tenant_id = var.tenant_id
  scanning_resource_group_id = var.scanning_resource_group_id
  storage_account_url = var.storage_account_url
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
  scan_multi_volume         = false
  scan_stopped_instances    = true
  query_text = var.query_text
  subscription_list = var.subscription_list
}

output "name" {
  value = lacework_integration_azure_agentless_scanning.example.name
}

output "client_id" {
  value = lacework_integration_azure_agentless_scanning.example.credentials[0].client_id
}

output "client_secret" {
  value = lacework_integration_azure_agentless_scanning.example.credentials[0].client_secret
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
