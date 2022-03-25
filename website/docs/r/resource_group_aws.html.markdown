---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group_aws"
description: |-
  Create and manage AWS Resource Groups
---

# lacework\_resource\_group\_aws

Use this resource to create an AWS Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://support.lacework.com/hc/en-us/articles/360041727354-Resource-Groups).

## Example Usage

```hcl
resource "lacework_resource_group_aws" "example" {
  name        = "My AWS Resource Group"
  description = "This groups a subset of AWS Accounts"
  accounts    = ["123456789", "234567891"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `accounts` - (Required) The list of AWS account ids to include in the resource group.
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework AWS Resource Group can be imported using a `RESOURCE_GUID`, e.g.

```
$ terraform import lacework_resource_group_aws.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `RESOURCE_GUID` from existing resource groups in your account, use the
Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
