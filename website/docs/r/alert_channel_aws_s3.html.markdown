---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_aws_s3"
description: |-
  Create and manage AWS S3 Alert Channel integrations
---

# lacework\_alert\_channel\_aws\_s3

Create the S3 alert channel that Lacework uses to send the data to your designated bucket in AWS

## Example Usage

```hcl
resource "lacework_alert_channel_aws_s3" "account_abc" {
  name = "s3 Alerts"
  bucket_arn  = "arn:aws:s3:::bucket_name/key_name"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `bucket_arn` - (Required) The ARN of the S3 bucket.
### Credentials

`credentials` supports the following arguments:

* `role_arn` - (Required) The ARN of the IAM role.
* `external_id` - (Required) The external ID for the IAM role.

## Import

A Lacework Aws S3 Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_aws_s3.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
