provider "lacework" {}

resource "lacework_integration_aws_cfg" "example" {
  name = "AWS config integration example"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
