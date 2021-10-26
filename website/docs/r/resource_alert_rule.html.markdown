---
subcategory: "Alert Rules"
layout: "lacework"
page_title: "Lacework: lacework_alert_rule"
description: |-
  Create and manage Lacework Alert Rules
---

# lacework\_alert\_rule

Use this resource to create a Lacework Alert Rule in order to categorize Lacework-identifiable assets.
For more information, see the [Alert Rules documentation](https://support.lacework.com/hc/en-us/articles/360042236733-Alert-Rules).

## Example Usage

```hcl
resource "lacework_alert_rule" "example" {
  name             = "My Alert Rule"
  description      = "This is an example alert rule"
  channels         = ["TECHALLY_AB90D4E77C93A9DE0DF6B22B9B06B9934645D6027C9D350"]
  severities       = ["High", "Medium"]
  event_categories = ["Compliance"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The alert rule name.
* `channels` - (Required) The list of alert channels for the rule to use.
* `severities` - (Required) The list of the severities that the rule will apply. `Critical`, `High`, `Medium`, `Low`, `Info`.
* `description` - (Optional) The description of the alert rule.
* `event_categories` - (Optional) The list of event categories the rule will apply to. `Compliance`, `App`, `Cloud`, 
  `File`, `Machine`, `User`, `Platform`.
* `resource_groups` - (Optional) The list of resource groups the rule will apply to.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Alert Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_alert_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
