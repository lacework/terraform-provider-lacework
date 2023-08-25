---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_pagerduty"
description: |-
  Create and manage PagerDuty Alert Channel integrations
---

# lacework\_alert\_channel\_pagerduty

Configure Lacework to forward alerts to PagerDuty through an API integration key.

## PagerDuty + Lacework Integration Benefits
* Extend Lacework Events to route to the correct people, at the correct time that fits your existing business processes using PagerDuty triage, escalations, and workflows.
* One-way event notification forwards from Lacework to PagerDuty.
* Lacework Alert Routing and Alert Rules settings allow you to configure which events and severities to receive and which resource groups and event categories you want events for. They grant complete control of the alert channels forwarded to PagerDuty.

## How it Works
Lacework events that arise from anomaly detection, compliance, vulnerabilities, or configured rule definitions
send an event to a [service](https://support.pagerduty.com/docs/services-and-integrations#section-configuring-services-and-integrations)
in PagerDuty. Events from Lacework can either trigger a new incident on the corresponding PagerDuty service or
be grouped as alerts into an existing incident.

For additional information about incidents and alerts, see https://support.pagerduty.com/docs/incidents and https://support.pagerduty.com/docs/alerts.

## Example Usage

```hcl
resource "lacework_alert_channel_pagerduty" "critical" {
  name            = "Forward Critical Alerts"
  integration_key = "1234abc8901abc567abc123abc78e012"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `integration_key` - (Required) The PagerDuty service integration key.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework PagerDuty Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_pagerduty.critical EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework alert-channel list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).

