provider "lacework" {}

data "lacework_agent_access_token" "k8s" {
  name = "k8s-deployments"
}

output "lacework_agent_access_token" {
  value     = data.lacework_agent_access_token.k8s.token
  sensitive = true
}
