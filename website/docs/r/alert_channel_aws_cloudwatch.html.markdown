---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_aws_cloudwatch"
description: |-
  Create and manage Amazon CloudWatch Alert Channel integrations
---

# lacework\_alert\_channel\_aws\_cloudwatch

Configure Lacework to forward alerts to an Amazon CloudWatch event bus.

-> **Note:** For more information about sending and receiving events between AWS accounts, refer to the Amazon [CloudWatch Events User Guide](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html).

## Example Usage

```hcl
resource "lacework_alert_channel_aws_cloudwatch" "all_events" {
  name            = "All events to default event-bus"
  event_bus_arn   = "arn:aws:events:us-west-2:1234567890:event-bus/default"
  group_issues_by = "Events"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `event_bus_arn` - (Required) The ARN of your AWS CloudWatch event bus.
* `group_issues_by` - (Optional) Defines how Lacework compliance events get grouped. Must be one of `Events` or `Resources`. Defaults to `Events`.
  The available options are:
  * **Events**:	Single AWS CloudWatch events will be created when compliance events of the same type but from different resources are detected by Lacework. For example, if three different S3 resources are generating the same compliance event, only one AWS event is created on the AWS CloudWatch event bus.
  * **Resources**: Multiple AWS CloudWatch events will be created when multiple resources are generating the same compliance event. For example, if three different S3 resources are generating the same compliance event, three AWS events are created on the AWS CloudWatch event bus.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Amazon CloudWatch Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_aws_cloudwatch.all_events EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).

