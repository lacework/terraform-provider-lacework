provider "lacework" {}

resource "lacework_integration_ecr" "iam_role" {
  name            = "ECR using IAM Role"
  registry_domain = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}

resource "lacework_integration_ecr" "access_key" {
  name            = "ECR using Access Keys"
  registry_domain = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  credentials {
    access_key_id     = "AWS123abcAccessKeyID"
    secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
  }

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["my-repo", "other-repo"]

  limit_by_labels = {
    key1 = "label1"
    key2 = "label2"
  }
}
