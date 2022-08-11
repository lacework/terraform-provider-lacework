---
subcategory: "Data Export Rules"
layout: "lacework"
page_title: "Lacework: lacework_data_export_rule"
description: |-
  Create and manage Lacework Data Export Rules
---

# lacework\_data\_export\_rule

Use this resource to export data collected from your Lacework account.
For more information, see the [Data Export Rules documentation](https://docs.lacework.com/console/category/data-shares--export).

## Example Usage

#### Data Export Rule with Slack Data Export Channel
```hcl
resource "lacework_data_export_rule" "example" {
  name             = "Data Export Rule From Terraform Updated"
  profile_versions = ["V1"]
  integration_ids  = ["INT_ABC123AB385C123D4567AB8EB45BA0E7ABCD12ABF65673A"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The data export rule name.
* `profile_versions` - (Required) The list of integration ids.
* `integration_ids` - (Required) The list of integration ids.
* `type` - (Optional) The type of the export rule. Defaults to `Dataexport`.
* `enabled` - (Optional) Whether the rule is enabled or disabled. Defaults to `true`.

## Import

A Lacework Data Export Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_data_export_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
