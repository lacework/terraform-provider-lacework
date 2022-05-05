---
subcategory: "Alert Profiles"
layout: "lacework"
page_title: "Lacework: lacework_alert_profile"
description: |-
  Create and manage Lacework Alert Profiles
---

# lacework\_alert\_profile

Use this resource to create a Lacework Alert Profile in order to map query results to events and form event descriptions.

## Example Usage

```hcl
resource "lacework_alert_profile" "example" {
  name    = "CUSTOM_PROFILE_TERRAFORM_TEST"
  extends = "LW_CFG_GCP_DEFAULT_PROFILE"

  alert {
    name        = "Violation"
    event_name  = "LW Configuration GCP Violation Alert"
    subject     = "{{_OCCURRENCE}} violation detected in project {{PROJECT_ID}}"
    description = var.alert_description
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The alert profile name, uniquely identifies the profile. Cannot start 'LW_' which is reserved for Lacework profiles.
* `extends` - (Required) The name of existing alert profile from which this profile extends.
* `alert` - (Required) The list of alert templates. See [Alert](#alert) below for details.
### Alert

`alert` supports the following arguments:

* `name` - (Required) The name that policies can use to refer to this template when generating alerts.
* `event_name` - (Required) The name of the resulting alert.
* `description` - (Required) The summary of the resulting alert.
* `subject` - (Required) A high-level observation of the resulting alert.

## Import

A Lacework Alert Profile can be imported using it's `name`, e.g.

```
$ terraform import lacework_alert_profile.example CUSTOM_PROFILE_TERRAFORM_TEST
```
