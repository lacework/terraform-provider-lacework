terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

data "lacework_agent_access_token" "k8s" {
  name = var.token_name
}

output "lacework_agent_access_token" {
  value     = data.lacework_agent_access_token.k8s.token
  sensitive = true
}

output "token_name" {
  value = data.lacework_agent_access_token.k8s.name
}

variable "token_name" {
  type    = string
  default = "k8s-deployments"
}
