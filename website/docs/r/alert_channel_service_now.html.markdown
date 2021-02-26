---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_service_now"
description: |-
  Create and manage Service Now Alert Channel integrations
---

# lacework\_alert\_channel\_service\_now

You can configure Lacework to forward alerts to Service Now using the ServiceNow REST API.

To find more information see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360005842314-ServiceNow).

## Example Usage

```hcl
resource "lacework_alert_channel_service_now" "example" {
  name         = "Service Now Alerts"
  instance_url = "snow-lacework.com"
  username     = "snow-user"
  password     = "snow-pass"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `instance_url` - (Required) The ServiceNow instance URL.
* `username` - (Required) The ServiceNow user name.
* `password` - (Required) The ServiceNow password.
* `custom_template_file` - (Optional) Populate fields in the ServiceNow incident with values from a custom template JSON file.
* `issue_grouping` - (Optional) Defines how Lacework compliance events get grouped. Must be one of `Events` or `Resources`. Defaults to `Events`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Service Now Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_service_now.ops_critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
