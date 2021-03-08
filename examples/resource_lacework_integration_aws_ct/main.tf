provider "lacework" {}

resource "lacework_integration_aws_ct" "example" {
  name      = "AWS CloudTrail integration example"
  queue_url = "https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
  retries = 10
}

resource "lacework_integration_aws_ct" "consolidated" {
  name      = "A consolidated CloudTrail example"
  queue_url = "https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }

  org_account_mappings {
    default_lacework_account = "lw_account_1"

    mapping {
      lacework_account = "lw_account_2"
      aws_accounts     = ["234556677", "774564564"]
    }

    mapping {
      lacework_account = "lw_account_3"
      aws_accounts     = ["553453453", "934534535"]
    }
  }
}
