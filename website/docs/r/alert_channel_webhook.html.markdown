---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_webhook"
description: |-
  Create and manage Webhook Alert Channel integrations
---

# lacework\_alert\_channel\_webhook

You can use this resource to create a custom webhook that receives Lacework alert notifications from a Lacework alert channel and forwards those alerts to a third-party application.

-> **Note:**  If the third-party application receiving the HTTP POST request requires an API token, enter the API token as part of the URL eg. https://webhook.com?api-token=123

To find more information about the alert payload, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360034367393-Webhook).

## Example Usage

```hcl
resource "lacework_alert_channel_webhook" "ops_critical" {
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
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
