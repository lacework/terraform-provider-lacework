---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_splunk"
description: |-
  Create and manage Splunk Alert Channel integrations
---

# lacework\_alert\_channel\_splunk

You can use this resource enabled Lacework to forward alerts to Splunk using an HTTP Event Collector.

To find more information see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360007889274-Splunk).

## Example Usage

```hcl
resource "lacework_alert_channel_" "ops_critical" {
  name      = "OPS Critical Alerts"
  hec_token = "AA111111-11AA-1AA1-11AA-11111AA1111A"
  host = "localhost"
  port = "80"
  event_data {
    index = "index"
    source = "source"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `channel` - (Optional) The Splunk channel name
* `hec_token` - (Required) The token you generate when you create a new HEC input.
* `host` - (Required) The hostname of the client from which you're sending data.
* `port` - (Required) The destination port for forwarding events [80 or 443].
* `ssl` - (Optional) Enable or Disable SSL.

* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Event Data

`event_data` supports the following arguments:

* `source` - (Required) The Splunk source.
* `index` - (Required) Index to store generated events.

## Import

A Lacework Splunk Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_splunk.ops_critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
