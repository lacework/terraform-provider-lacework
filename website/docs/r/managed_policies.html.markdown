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
    enabled  = false
    severity = "Critical"
  }
  policy {
    id       = "lacework-global-10"
    severity  = "Low"
  }
}
```

## Argument Reference

For each `policy` block, the following arguments are supported:

* `id` - (Required) The Lacework-defined policy id.
* `enabled` - (Required) Whether the policy is enabled or disabled.
* `severity` - (Required) The list of the severities. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.
