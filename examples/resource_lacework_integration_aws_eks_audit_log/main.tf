terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_integration_aws_eks_audit_log" "example" {
  name      = "AWS EKS audit log integration example"
  sns_arn = "arn:aws:sns:us-west-2:123456789123:foo-lacework-eks"
  credentials {
    role_arn    = "arn:aws:iam::249446771485:role/lacework-iam-example-role"
    external_id = "12345"
  }
  retries = 10
}
