terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_resource_group" "example" {
  type        = "AWS"
  name        = var.resource_group_name
  description = var.description

  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Region"
      operation = "EQUALS"
      value     = ["us-east-1"]
    }

    filter {
      filter_name = "filter2"
      field     = "Region"
      operation = "EQUALS"
      value     = ["us-west-2"]
    }

    group {
      operator = "AND"
      group {
        operator = "OR"
        filter {
          filter_name = "filter3"
          field     = "Account"
          operation = "EQUALS"
          value     = ["987654321"]
        }
        filter {
          filter_name = "filter4"
          field     = "Account"
          operation = "EQUALS"
          value     = ["123456789"]
        }
      }
    }
  }
}

