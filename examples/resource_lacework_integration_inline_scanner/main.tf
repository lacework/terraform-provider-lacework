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
  name                   = var.integration_name

  limit_num_scan        = 60
  identifier_tag {
    key   = "key"
    value = "value"
  }
}
