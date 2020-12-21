provider "lacework" {}

resource "lacework_alert_channel_aws_s3" "s3_export" {
  name        = "S3 Data Export"
  credentials {
    external_id = "12345"
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    bucket_arn  = "arn:aws:s3:::bucket_name/key_name"
  }
}