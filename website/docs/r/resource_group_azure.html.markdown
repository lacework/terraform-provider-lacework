---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group_azure"
description: |-
  Create and manage Azure Resource Groups
---

# lacework\_resource\_group\_azure

Use this resource to create an Azure Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://support.lacework.com/hc/en-us/articles/360041727354-Resource-Groups).

## Example Usage

```hcl
resource "lacework_resource_group_azure" "example" {
  name          = "My Azure Resource Group"
  description   = "This groups a subset of Azure Subscriptions"
  tenant        = "a11aa1ab-111a-11ab-a000-11aa1111a11a"
  subscriptions = ["1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0", "2b000c3-ab10-1a01-1abc-1a000ab0a0a0"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `tenant` - (Required) The Azure tenant id.
* `subscriptions` - (Required) The list of Azure subscription ids to include in the resource group.
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Azure Resource Group can be imported using a `RESOURCE_GUID`, e.g.

```
$ terraform import lacework_resource_group_azure.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `RESOURCE_GUID` from existing resource groups in your account, use the
Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
