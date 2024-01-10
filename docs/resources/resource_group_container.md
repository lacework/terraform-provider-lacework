---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group_container"
description: |-
  Create and manage Container Resource Groups
---

# lacework\_resource\_group\_container

Use this resource to create a Container Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://docs.lacework.net/console/resource-groups).

## Example Usage

```hcl
resource "lacework_resource_group_container" "example" {
  name           = "My Container Resource Group"
  description    = "This groups a subset of Container Tags"
  container_tags = ["my-container"]
  container_label {
    key   = "name"
    value = "my-container"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `container_labels` - (Required) The key value pairs of container labels to include in the resource group.
* `container_tags` - (Required) The list of container tags to include in the resource group.
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Container Resource Group can be imported using a `RESOURCE_GUID`, e.g.

```
$ terraform import lacework_resource_group_container.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `RESOURCE_GUID` from existing resource groups in your account, use the
Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
