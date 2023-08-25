---
subcategory: "Data Export Rules"
layout: "lacework"
page_title: "Lacework: lacework_data_export_rule"
description: |-
  Create and manage Lacework Data Export Rules
---

# lacework\_data\_export\_rule

Use this resource to export data collected from your Lacework account and send it to an S3 bucket of your choice.

For more information, see [Data Export Rules](https://docs.lacework.com/console/category/data-shares--export) and
[S3 Data Export](https://docs.lacework.com/console/s3-data-export) documentation.

## Example Usage

#### Data Export Rule
```hcl
resource "lacework_data_export_rule" "example" {
  name            = "Data Export Rule From Terraform Updated"
  integration_ids = ["INT_ABC123AB385C123D4567AB8EB45BA0E7ABCD12ABF65673A"]
}
```

#### Data Export Rule with S3 Data Export Channel
```hcl
resource "lacework_alert_channel_aws_s3" "data_export" {
  name       = "s3 data export to account 1234567890"
  bucket_arn = "arn:aws:s3:::bucket_name/key_name"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}

resource "lacework_data_export_rule" "example" {
  name            = "business unit export rule"
  integration_ids = [lacework_alert_channel_aws_s3.data_export.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The data export rule name.
* `integration_ids` - (Required) The list s3 data export alert channel ids for the rule to use.
* `description` - (Optional) The summary of the data export rule.
* `enabled` - (Optional) Whether the rule is enabled or disabled. Defaults to `true`.

## Import

A Lacework Data Export Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_data_export_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
