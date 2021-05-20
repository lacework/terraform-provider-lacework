provider "lacework" {}

resource "lacework_integration_aws_govcloud_cfg" "example" {
  name = "AWS gov cloud config integration example"
  account_id = "553453453"
  credentials {
    access_key_id     = "AWS123abcAccessKeyID"
    secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
  }

  retries = 10
}
