---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_email"
description: |-
  Create and manage Email Alert Channel integrations
---

# lacework\_alert\_channel\_email

Lacework can generate and send alert summaries and reports to email addresses using an email alert channel. By default,
Lacework creates a single email alert channel during the initial Lacework onboarding process and new members are added
automatically. The default channel cannot be edited. You can add more email alert channels.

To find more information about the Email alert channel integration, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360023638654-Email).

## Example Usage

```hcl
resource "lacework_alert_channel_email" "auditors" {
  name       = "Auditors Alerts"
  recipients = [
    "my@example.com",
    "alias@example.com"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `recipients` - (Required) The list of email addresses that you want to receive the alerts.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Email Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_email.auditors EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework alert-channel list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
