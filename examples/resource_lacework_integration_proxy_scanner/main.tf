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
  name                   = var.name

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["repo/my-image", "repo/other-image"]

  limit_by_label {
    key   = "foo"
    value = "bar"
  }
}

output "server_token" {
    value = lacework_integration_proxy_scanner.example.server_token
}