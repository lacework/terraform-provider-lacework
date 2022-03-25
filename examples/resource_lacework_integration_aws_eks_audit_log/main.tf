provider "lacework" {}

resource "lacework_integration_aws_eks_audit_log" "example" {
  name      = "AWS EKS audit log integration example"
  sns_arn = "arn:aws:sns:us-west-2:123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
  retries = 10
}
