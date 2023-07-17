---
subcategory: "Policies"
layout: "lacework"
page_title: "Lacework: lacework_managed_policies"
description: |-
  Manage Lacework Defined Policies
---

# lacework\_managed\_policies

Use this resource to update the "enabled" and the "severity" property for Lacework defined policies.

## Example Usage

Create a lacework_managed_policies resource to manage three Lacework defined policies.

```hcl
resource "lacework_managed_policies" "example" {
  policy {
    id       = "lacework-global-1"
    enabled  = true
    severity = "high"
  }
  policy {
    id       = "lacework-global-2"
    enabled  = true
    severity = "critical"
  }
  policy {
    id       = "lacework-global-10"
    enabled  = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The Lacework defined policy id.
* `enabled` - (Required) Whether the policy is enabled or disabled.
* `severity` - (Optional) The list of the severities. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.
