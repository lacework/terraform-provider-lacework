terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_ecr" "iam_role" {
  name            = var.integration_name
  registry_domain = var.registry_domain
  credentials {
    role_arn    = var.role_arn
    external_id = var.external_id
  }
  non_os_package_support = var.non_os_packages
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
