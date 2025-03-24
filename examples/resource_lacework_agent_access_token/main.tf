terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_agent_access_token" "k8s" {
  name        = var.token_name
  description = "Token for K8S clusters"
  os          = var.os_type
}

variable "token_name" {
  type    = string
  default = "k8s-deployments"
}

variable "os_type" {
  type    = string
  default = "linux"
}

output "token_name" {
  value = lacework_agent_access_token.k8s.name
}