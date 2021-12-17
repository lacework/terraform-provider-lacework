---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_newrelic"
description: |-
  Create and manage New Relic Insights Alert Channel integrations
---

# lacework\_alert\_channel\_newrelic

You can configure a Lacework alert channel to forward alerts to New Relic using the Insights API.

To find more information about the New Relic Insights alert channel integration, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360005842354-New-Relic).

## Example Usage

```hcl
resource "lacework_alert_channel_newrelic" "example" {
  name       = "Example New Relic Insights Alert"
  account_id = 2338053
  insert_key = "x-xx-xxxxxxxxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `account_id` - (Required) The New Relic account ID.
* `insert_key` - (Required) The New Relic Insert API key.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework New Relic Insights Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_newrelic.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
