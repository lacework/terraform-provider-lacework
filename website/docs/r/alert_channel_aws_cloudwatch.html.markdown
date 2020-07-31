---
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_aws_cloudwatch"
description: |-
  Create an manage AWS CloudWatch Alert Channel integrations
---

# lacework\_alert\_channel\_aws\_cloudwatch

Configure Lacework to forward alerts to an AWS CloudWatch event bus.

-> **Note:** For more information about sending and receiving events between AWS accounts, refer to the Amazon [CloudWatch Events User Guide](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html).

## Example Usage

```hcl
resource "lacework_alert_channel_aws_cloudwatch" "all_events" {
  name               = "All events to default event-bus"
  event_bus_arn      = "arn:aws:events:us-west-2:1234567890:event-bus/default"
  min_alert_severity = 5
  group_issues_by    = "EVENTS"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `event_bus_arn` - (Required) The ARN of your AWS CloudWatch event bus.
* `min_alert_severity` - (Optional) The minimum severity level to alert. Defaults to `3`.
  The available levels are:
  * **1**: Critical Alerts only
  * **2**: High Alerts and above
  * **3**: Medium Alerts and above
  * **4**: Low Alerts and above
  * **5**: All Alerts
* `group_issues_by` - (Optional) Defines how Lacework compliance events get grouped. Must be one of `EVENTS` or `RESOURCES`. Defaults to `EVENTS`.
  The available options are:
  * **EVENTS**:	Single AWS CloudWatch events will be created when compliance events of the same type but from different resources are detected by Lacework. For example, if three different S3 resources are generating the same compliance event, only one AWS event is created on the AWS CloudWatch event bus.
  * **RESOURCES**: Multiple AWS CloudWatch events will be created when multiple resources are generating the same compliance event. For example, if three different S3 resources are generating the same compliance event, three AWS events are created on the AWS CloudWatch event bus.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework AWS CloudWatch Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_aws_cloudwatch.all_resources EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).

