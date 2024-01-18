terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

data "lacework_metric_module" "test" {
  name    = "terraform-aws-cloudtrail"
  version = "1.0.0"
}
