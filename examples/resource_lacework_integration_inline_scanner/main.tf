terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {
  profile="snifftest-composite"
}

resource "lacework_integration_inline_scanner" "example" {
  name                   = var.name

  limit_num_scan        = 60
  identifier_tag {
    key   = "foo"
    value = "bar"
  }
}

output "server_token" {
    value = lacework_integration_inline_scanner.example.server_token
}