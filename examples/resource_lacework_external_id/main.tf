terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

# Format
#
#    lweid:<csp>:<version>:<lw_tenant_name>:<aws_acct_id>:<random_string_size_10>
#
resource "lacework_external_id" "aws_123456789012" {
  csp        = "aws"
  account_id = "123456789012"
}

# Example output
#
#    lweid:aws:v2:customerdemo:123456789012:dkl31.09ip
#
output "external_id" {
  value = lacework_external_id.aws_123456789012.v2
}

