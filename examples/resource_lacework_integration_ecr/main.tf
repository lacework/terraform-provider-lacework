provider "lacework" {}

resource "lacework_integration_ecr" "example" {
  name              = "ERC Example"
  registry_domain   = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  access_key_id     = "AWS123abcAccessKeyID"
  secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
  limit_by_tag      = "dev*"
  limit_by_label    = "*label"
  limit_by_repos    = "my-repo,other-repo"
  limit_num_imgs    = 10
}
