---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_victorops"
description: |-
  Create and manage VictorOps Alert Channel integrations
---

# lacework\_alert\_channel\_victorops

You can configure Lacework to forward alerts to specific VictorOps groups using a VictorOps REST endpoint.

To find more information about the alert payload, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360005916533-VictorOps).

## Example Usage

```hcl
resource "lacework_alert_channel_victorops" "example" {
  name       = "VictorOps Example"
  webhook_url = "https://alert.victorops.com/integrations/generic/20131114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `webhook_url` - (Required) The URL of your VictorOps webhook that will receive the HTTP POST.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework VictorOps Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_victorops.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
