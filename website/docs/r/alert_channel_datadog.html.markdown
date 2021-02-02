---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_datadog"
description: |-
  Create and manage Datadog Alert Channel integrations
---

# lacework\_alert\_channel\_datadog

The Datadog alert channel provides a unified view of your metrics, logs, and performance data combined with your cloud security data.

To find more information about the Datadog alert integration, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360036989953-Datadog).

## Example Usage

```hcl
resource "lacework_alert_channel_" "ops_critical" {
  name      = "OPS Critical Alerts"
  datadog_site = "eu"
  datadog_service = "Events Summary"
  api_key = "datadog-key"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `datadog_site` - (Required) Where to store your logs, either the US or Europe. Must be one of `com` or `eu`. Defaults to `com`.
* `datadog_service` - (Required) `Logs Detail`, `Logs Summary`, or `Events Summary`. Defaults to `Logs Detail`.
* `api_key` - (Required) The URL of your datadog that will receive the HTTP POST.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Datadog Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_datadog.ops_critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
