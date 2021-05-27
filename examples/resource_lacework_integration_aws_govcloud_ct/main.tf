provider "lacework" {}

resource "lacework_integration_aws_govcloud_ct" "example" {
  name       = "AWS gov cloud config integration example"
  account_id = "553453453"
  queue_url  = "https://sqs.us-gov-west-1.amazonaws.com/123456789012/my_queue"
  credentials {
    access_key_id     = "AWS123abcAccessKeyID"
    secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
  }

  retries = 10
}
