---
subcategory: "Policies"
layout: "lacework"
page_title: "Lacework: lacework_managed_policies"
description: |-
  Manage Lacework Defined Policies
---

# lacework\_managed\_policies

Use this resource to update the `state` (enabled/disabled) and the `severity` properties for Lacework-defined policies.

## Example Usage

The following example shows how to manage three Lacework-defined policies.

```hcl
resource "lacework_managed_policies" "example" {
  policy {
    id       = "lacework-global-1"
    enabled  = true
    severity = "High"
  }
  policy {
    id       = "lacework-global-2"
    enabled  = true
    severity = "Critical"
  }
  policy {
    id       = "lacework-global-10"
    enabled  = false
  }
}
```

## Argument Reference

For each `policy` block, the following arguments are supported:

* `id` - (Required) The Lacework-defined policy id.
* `enabled` - (Required) Whether the policy is enabled or disabled.
* `severity` - (Optional) The list of the severities. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.

## Import

A lacework_managed_policies resource can be imported using a UTC string as the ID, e.g.

```
$ terraform import lacework_managed_policies.example "2023-07-19 20:58:07.320676 +0000 UTC"
```
