---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_qradar"
description: |-
  Create and manage IBM QRadar Alert Channel integrations
---

# lacework\_alert\_channel\_qradar

You can configure an alert channel to send Lacework alert notifications to IBM QRadar.

To find more information about the IBM QRadar alert channel integration, see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360056898693-IBM-QRadar).

## Example Usage

```hcl
resource "lacework_alert_channel_qradar" "example" {
  name               = "QRadar Channel Alert Example"
  host_url           = "https://qradar-lacework.com"
  host_port          = 4000
  communication_type = "HTTPS"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `host_url` - (Required) The domain name or IP address of QRadar.
* `host_port` - (Required) The listen port defined in QRadar.
* `communication_type` - (Required) The communication protocol used. Must be one of `HTTPS` or `HTTPS Self Signed Cert`. 
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework IBM QRadar Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_qradar.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
