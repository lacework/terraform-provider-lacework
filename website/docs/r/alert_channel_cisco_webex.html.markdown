---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_cisco_webex"
description: |-
  Create and manage Cisco Webex Alert Channel integrations
---

# lacework\_alert\_channel\_cisco\_webex

You can configure a Lacework alert channel to forward alerts to a Cisco Webex Teams space as an incoming webhook.

To find more information about the Cisco Webex alert channel integration, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360005840154-Cisco-Webex-Teams).

## Example Usage

```hcl
resource "lacework_alert_channel_cisco_webex" "example" {
  name       = "Example Cisco Webex Alert"
  webhook_url = "https://webexapis.com/v1/webhooks/incoming/api-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `webhook_url` - (Required) The Cisco Webex webhook URL.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Cisco Webex Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_cisco_webex.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework alert-channel list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
