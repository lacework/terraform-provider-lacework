terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_integration_aws_agentless_scanning" "example" {
  name                      = var.name
  query_text                = var.query_text
  scan_frequency            = 24
  scan_containers           = true
  scan_host_vulnerabilities = true
}

variable "name" {
  type    = string
  default = "AWS Agentless Scanning Example"
}

variable "query_text" {
  type    = string
  default = ""
}

output "name" {
  value = lacework_integration_aws_agentless_scanning.example.name
}