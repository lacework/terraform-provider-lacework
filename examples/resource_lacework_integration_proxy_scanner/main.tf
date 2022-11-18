terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {
  profile = "snifftest-composite"
}

resource "lacework_integration_proxy_scanner" "example" {
  name                   = var.integration_name

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["repo/my-image", "repo/other-image"]

  limit_by_label {
    key   = "foo"
    value = "bar"
  }

  policy_evaluate = false
  policy_guids = ["VULN_0595430C23E5C3BBB5EBDB59CEF17467AF592C825562090FDA9"]
}

output "server_token" {
    value = lacework_integration_proxy_scanner.example.server_token
}

output "server_uri" {
    value = lacework_integration_proxy_scanner.example.server_uri
}

output "policy_evaluate" {
    value = lacework_integration_proxy_scanner.example.policy_evaluate
}

output "policy_guids" {
    value = lacework_integration_proxy_scanner.example.policy_guids
}