---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_microsoft_teams"
description: |-
  Create and manage Microsoft Teams Alert Channel integrations
---

# lacework\_alert\_channel\_microsoft_teams

You can configure Lacework to forward alerts to a Microsoft Teams channel through an incoming webhook.
[Create a incoming webhook](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook)

To find more information about the alert payload, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360051656053-Microsoft-Teams).

## Example Usage

```hcl
resource "lacework_alert_channel_microsoft_teams" "example" {
  name      = "Microsoft Teams Alerts"
  webhook_url = "https://outlook.office.com/webhook/api-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `webhook_url` - (Required) The URL of your Microsoft Teams incoming webhook.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Webhook Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_microsoft_teams.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
