terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

data "lacework_metrics" "test" {
  name    = "terraform-aws-cloudtrail"
  version = "1.0.0"
}

output "lacework_trace_id" {
  value = data.lacework_metrics.test.trace_id
}
