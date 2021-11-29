---
subcategory: "Report Rules"
layout: "lacework"
page_title: "Lacework: lacework_report_rule"
description: |-
  Create and manage Lacework Report Rules
---

# lacework\_report\_rule

Use this resource to create a Lacework Report Rule in order to route reports to an email alert channel.
For more information, see the [Report Rules documentation](https://docs.lacework.com/report-rules).

## Example Usage

#### Report Rule with Aws Compliance Reports
```hcl
resource "lacework_report_channel_email" "team_email" {
  name       = "Team Emails"
  recipients = ["foo@example.com", "bar@example.com"]
}

resource "lacework_report_rule" "example" {
  name                 = "My Report Rule"
  description          = "This is an example report rule"
  email_alert_channels = [lacework_report_channel_email.team_email.id]
  severities           = ["Critical", "High"]

  aws_compliance_reports {
    cis_s3 = true
  }
  weekly_snapshot = true
}
```

#### Report Rule with Gcp Compliance Reports and Gcp Resource Group
```hcl
resource "lacework_report_channel_email" "team_email" {
  name       = "Team Emails"
  recipients = ["foo@example.com", "bar@example.com"]
}

resource "lacework_resource_group_gcp" "all_gcp_projects" {
  name         = "GCP Resource Group"
  description  = "All Gcp Projects"
  organization = "MyGcpOrg"
  projects     = ["*"]
}

resource "lacework_report_rule" "example" {
  name             = "My Report Rule"
  description      = "This is an example report rule"
  channels         = [lacework_report_channel_slack.ops_critical.id]
  resource_groups  = [lacework_resource_group_gcp.all_gcp_projects.id]
  severities       = ["Critical"]

  gcp_compliance_reports {
    k8s = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The report rule name.
* `channels` - (Required) The list of report channels for the rule to use.
* `severities` - (Required) The list of the severities that the rule will apply. Valid severities include: 
  `Critical`, `High`, `Medium`, `Low` and `Info`.
* `description` - (Optional) The description of the report rule.
* `resource_groups` - (Optional) The list of resource groups the rule will apply to.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `aws_compliance_reports` - (Optional) Compliance reports for Aws. See [Aws Compliance Reports](#Aws Compliance Reports) below for details.
* `azure_compliance_reports` - (Optional) Compliance reports for Azure. See [Azure Compliance Reports](#Azure Compliance Reports) below for details.
* `gcp_compliance_reports` - (Optional) Compliance reports for Gcp. See [Gcp Compliance Reports](#Gcp Compliance Reports) below for details.
* `daily_compliance_reports` - (Optional) Daily event summary reports. See [Daily Compliance Reports](#Daily Compliance Reports) below for details.
* `weekly_snapshot` - (Optional) A weekly compliance trend report for all monitored resources. Defaults to `false`.

### Aws Compliance Reports

`aws_compliance_reports` supports the following arguments:

* `cis_s3` - (Optional) Defaults to `false`.
* `hipaa` - (Optional) Defaults to `false`.
* `iso_2700` - (Optional) Defaults to `false`.
* `nist_800_53_rev4` - (Optional) Defaults to `false`.
* `nist_800_171_rev2` - (Optional) Defaults to `false`.
* `pci` - (Optional) Defaults to `false`.
* `soc` - (Optional) Defaults to `false`.
* `soc_rev2` - (Optional) Defaults to `false`.

### Azure Compliance Reports

`azure_compliance_reports` supports the following arguments:

* `cis` - (Optional) Defaults to `false`.
* `cis_131` - (Optional) Defaults to `false`.
* `pci` - (Optional) Defaults to `false`.
* `soc` - (Optional) Defaults to `false`.

### Gcp Compliance Reports

`gcp_compliance_reports` supports the following arguments:

* `cis` - (Optional) Defaults to `false`.
* `hipaa` - (Optional) Defaults to `false`.
* `hipaa_rev2` - (Optional) Defaults to `false`.
* `iso_27001` - (Optional) Defaults to `false`.
* `cis_12` - (Optional) Defaults to `false`.
* `k8s` - (Optional) Defaults to `false`.
* `pci` - (Optional) Defaults to `false`.
* `pci_rev2` - (Optional) Defaults to `false`.
* `soc` - (Optional) Defaults to `false`.
* `soc_rev2` - (Optional) Defaults to `false`.

### Daily Compliance Reports

`daily_compliance_reports` supports the following arguments:

* `host_security` - (Optional) Defaults to `false`.
* `openshift_compliance` - (Optional) Defaults to `false`.
* `openshift_compliance_events` - (Optional) Defaults to `false`.
* `host_security` - (Optional) Defaults to `false`.
* `aws_cloudtrail` - (Optional) Defaults to `false`.
* `host_security` - (Optional) Defaults to `false`.
* `aws_compliance` - (Optional) Defaults to `false`.
* `azure_activity_log` - (Optional) Defaults to `false`.
* `gcp_audit_trail` - (Optional) Defaults to `false`.
* `gcp_compliance` - (Optional) Defaults to `false`.

## Import

A Lacework Report Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_report_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
