---
subcategory: "Alert Rules"
layout: "lacework"
page_title: "Lacework: lacework_alert_rule"
description: |-
  Create and manage Lacework Alert Rules
---

# lacework\_alert\_rule

Use this resource to create a Lacework Alert Rule in order to route events to the appropriate people or tools.
For more information, see the [Alert Rules documentation](https://support.lacework.com/hc/en-us/articles/360042236733-Alert-Rules).

## Example Usage

#### Alert Rule with Slack Alert Channel
```hcl
resource "lacework_alert_channel_slack" "ops_critical" {
  name      = "OPS Critical Alerts"
  slack_url = "https://hooks.slack.com/services/ABCD/12345/abcd1234"
}

resource "lacework_alert_rule" "example" {
  name                = "My Alert Rule"
  description         = "This is an example alert rule"
  alert_channels      = [lacework_alert_channel_slack.ops_critical.id]
  severities          = ["Critical"]
  alert_subcategories = ["Compliance"]
  alert_categories    = ["Policy"]
  alert_sources       = ["AWS"]
}
```

#### Alert Rule with Slack Alert Channel and Gcp Resource Group
```hcl
resource "lacework_alert_channel_slack" "ops_critical" {
  name      = "OPS Critical Alerts"
  slack_url = "https://hooks.slack.com/services/ABCD/12345/abcd1234"
}

resource "lacework_resource_group_gcp" "all_gcp_projects" {
  name         = "GCP Resource Group"
  description  = "All Gcp Projects"
  organization = "MyGcpOrg"
  projects     = ["*"]
}

resource "lacework_alert_rule" "example" {
  name                = "My Alert Rule"
  description         = "This is an example alert rule"
  alert_channels      = [lacework_alert_channel_slack.ops_critical.id]
  severities          = ["Critical"]
  alert_subcategories = ["Compliance"]
  resource_groups     = [lacework_resource_group_gcp.all_gcp_projects.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The alert rule name.
* `alert_channels` - (Required) The list of alert channels for the rule to use.
* `severities` - (Required) The list of the severities that the rule will apply. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.
* `description` - (Optional) The description of the alert rule.
* `alert_subcategories` - (Optional) The list of alert subcategories the rule will apply to. Valid categories include:
  `Compliance`, `Application`, `Cloud Activity`, `File`, `Machine`, `User`, `Platform`, `Kubernetes Activity`, `Registry`, `SystemCall`,`Host Vulnerability`, `Container Vulnerability`, `Threat Intel`.
* `alert_categories` - (Optional) The alert categories that will use this rule for alert routing. Valid categories include:
  `Anomaly`, `Policy`, `Composite`.
* `alert_sources` - (Optional) The alert sources that will use this rule for alert routing. Valid sources include:
  `Agent`, `AWS`, `Azure`, `GCP`, `K8s`.
* `resource_groups` - (Optional) The list of resource groups the rule will apply to.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `event_categories` - (Optional, **Deprecated**) The list of event categories the rule will apply to. Valid categories include:
    `Compliance`, `App`, `Cloud`, `File`, `Machine`, `User`, `Platform`, `K8sActivity`, `Registry`, `SystemCall`.
This attribute is deprecated use `alert_subcategories` instead.



## Import

A Lacework Alert Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_alert_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
