provider "lacework" {}

resource "lacework_integration_aws_ct" "example" {
  name      = "AWS CloudTrail integration example"
  queue_url = "https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
