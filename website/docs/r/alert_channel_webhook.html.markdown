---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_webhook"
description: |-
  Create and manage Webhook Alert Channel integrations
---

# lacework\_alert\_channel\_webhook

Configure Lacework to forward alerts to a 3rd party Webhook.

-> **Note:**  If the third-party that is receiving the HTTP POST request requires an API token, enter the API Token as part of the URL eg. https://webhook.com?api-token=123

## Example Usage

```hcl
resource "lacework_alert_channel_" "ops_critical" {
  name      = "OPS Critical Alerts"
  webhook_url = "https://webhook.com?api-token=123"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `webhook_url` - (Required) The URL of your webhook that will receive the HTTP POST.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Webhook Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_webhook.ops_critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
