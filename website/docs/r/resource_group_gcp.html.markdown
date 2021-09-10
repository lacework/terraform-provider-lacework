---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group_gcp"
description: |-
Create and manage GCP Resource Groups
---

# lacework\_resource\_group\_gcp

Use this resource to create a GCP Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://support.lacework.com/hc/en-us/articles/360041727354-Resource-Groups).

## Example Usage

```hcl
resource "lacework_resource_group_gcp" "example" {
  name         = "My GCP Resource Group"
  description  = "This groups a subset of Gcp Projects"
  projects     = ["project-1", "project-2", "project-3"]
  organization = "MyGcpOrgID"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `projects` - (Required) The list of GCP project IDs to include in the resource group.
* `organization` - (Required) The GCP organization ID.
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework GCP Resource Group can be imported using a `RESOURCE_GUID`, e.g.

```
$ terraform import lacework_resource_group_gcp.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `RESOURCE_GUID` from existing resource groups in your account, use the
Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
