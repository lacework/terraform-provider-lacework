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
}

variable "token_name" {
  type    = string
  default = "k8s-deployments"
}

output "token_name" {
  value = lacework_agent_access_token.k8s.name
}