---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group_account"
description: |-
  Create and manage Lacework Account Resource Groups
---

# lacework\_resource\_group\_account

Use this resource to create a Lacework Account Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://support.lacework.com/hc/en-us/articles/360041727354-Resource-Groups).

## Example Usage

```hcl
resource "lacework_resource_group_account" "example" {
  name        = "My Lacework Account Resource Group"
  description = "This groups a subset of Lacework accounts"
  accounts    = ["my-account"]
}
```

### Organization Level Access

The `lacework_resource_group_account` is an organization-level resource that requires the Lacework provider to be configured with the `organization` argument to `true`.

```hcl
provider "lacework" {
  organization = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `accounts` - (Required) The list of Lacework accounts to include in the resource group.
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Account Resource Group can be imported using a `RESOURCE_GUID`, e.g.

```
$ terraform import lacework_resource_group_account.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `RESOURCE_GUID` from existing resource groups in your account, use the
Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
