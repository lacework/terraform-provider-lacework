---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_slack"
description: |-
  Create and manage Slack Alert Channel integrations
---

# lacework\_alert\_channel\_slack

Configure Lacework to forward alerts to a Slack channel through an incoming webhook.

-> **Note:** Lacework recommends creating a dedicated Slack channel for Lacework events.

## Example Usage

```hcl
resource "lacework_alert_channel_slack" "ops_critical" {
  name      = "OPS Critical Alerts"
  slack_url = "https://hooks.slack.com/services/ABCD/12345/abcd1234"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `slack_url` - (Required) The URL of the incoming Slack webhook.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Slack Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_slack.ops_critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
