---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_aws_s3"
description: |-
  Create and manage AWS S3 Alert Channel integrations
---

# lacework\_alert\_channel\_aws\_s3

S3 data export allows you to export data collected from your Lacework account and send it to an S3 bucket of
your choice. You can extend Lacework processed/normalized data to report/visualize alone or combine with other
business/security data to get insights and make meaningful business decisions.

!> **Warning:** This feature is currently in beta.

Every hour, Lacework collects data from your Lacework account and sends it to an internal Lacework S3 bucket as
a staging location. The data remains in the internal Lacework S3 bucket until its hourly scheduled export to your
designated S3 bucket.

For detailed information about the data exported by Lacework, see [Lacework Data Share](https://support.lacework.com/hc/sections/360011719393).

-> **Note:** Before proceeding, ensure that the bucket that will receive the data from Lacework already exists in AWS.

## Example Usage

```hcl
resource "lacework_alert_channel_aws_s3" "data_export" {
  name = "s3 Alerts"
  bucket_arn  = "arn:aws:s3:::bucket_name/key_name"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `bucket_arn` - (Required) The ARN of the S3 bucket.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `role_arn` - (Required) The ARN of the IAM role.
* `external_id` - (Required) The external ID for the IAM role.

## Import

A Lacework AWS S3 Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_aws_s3.data_export EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
