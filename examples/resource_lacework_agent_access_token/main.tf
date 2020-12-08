provider "lacework" {}

resource "lacework_agent_access_token" "k8s" {
  name        = "k8s-deployments"
  description = "Token for K8S clusters"
}
