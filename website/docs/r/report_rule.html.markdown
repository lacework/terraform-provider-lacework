---
subcategory: "Report Rules"
layout: "lacework"
page_title: "Lacework: lacework_report_rule"
description: |-
  Create and manage Lacework Report Rules
---

# lacework\_report\_rule

Use this resource to create a Lacework Report Rule in order to route reports to one or more email alert channels.
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
  name                 = "My Report Rule"
  description          = "This is an example report rule"
  email_alert_channels = [lacework_report_channel_email.team_email.id]
  resource_groups      = [lacework_resource_group_gcp.all_gcp_projects.id]
  severities           = ["Critical"]

  gcp_compliance_reports {
    k8s = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The report rule name.
* `email_alert_channels` - (Required) The list of email alert channels for the rule to use.
* `severities` - (Required) The list of the severities that the rule will apply. Valid severities include: 
  `Critical`, `High`, `Medium`, `Low` and `Info`.
* `description` - (Optional) The description of the report rule.
* `resource_groups` - (Optional) The list of resource groups the rule will apply to.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `aws_compliance_reports` - (Optional) Compliance reports for Aws. See [Aws Compliance Reports](#aws-compliance-reports) below for details.
* `azure_compliance_reports` - (Optional) Compliance reports for Azure. See [Azure Compliance Reports](#azure-compliance-reports) below for details.
* `gcp_compliance_reports` - (Optional) Compliance reports for Gcp. See [Gcp Compliance Reports](#gcp-compliance-reports) below for details.
* `daily_compliance_reports` - (Optional) Daily event summary reports. See [Daily Compliance Reports](#faily-compliance-reports) below for details.
* `weekly_snapshot` - (Optional) A weekly compliance trend report for all monitored resources. Defaults to `false`.

### Aws Compliance Reports

`aws_compliance_reports` supports the following arguments:

* `cis_s3` - (Optional) AWS CIS Benchmark and S3 Report. Defaults to `false`.
* `hipaa` - (Optional) AWS HIPAA Report. Defaults to `false`.
* `iso_2700` - (Optional) AWS ISO 27001:2013 Report. Defaults to `false`.
* `nist_800_53_rev4` - (Optional) AWS NIST 800-53 Report. Defaults to `false`.
* `nist_800_171_rev2` - (Optional) AWS NIST 800-171 Report. Defaults to `false`.
* `pci` - (Optional) AWS PCI DSS Report. Defaults to `false`.
* `soc` - (Optional) AWS SOC 2 Report. Defaults to `false`.
* `soc_rev2` - (Optional) AWS SOC 2 Report Rev2. Defaults to `false`.

### Azure Compliance Reports

`azure_compliance_reports` supports the following arguments:

* `cis` - (Optional) Azure CIS Benchmark. Defaults to `false`.
* `cis_131` - (Optional) Azure CIS 1.3.1 Benchmark. Defaults to `false`.
* `pci` - (Optional) Azure PCI Benchmark. Defaults to `false`.
* `soc` - (Optional) Azure SOC 2 Report. Defaults to `false`.

### Gcp Compliance Reports

`gcp_compliance_reports` supports the following arguments:

* `cis` - (Optional) GCP CIS Benchmark. Defaults to `false`.
* `hipaa` - (Optional) GCP HIPAA Report. Defaults to `false`.
* `hipaa_rev2` - (Optional) GCP HIPAA Report Rev2. Defaults to `false`.
* `iso_27001` - (Optional) GCP ISO 27001 Report. Defaults to `false`.
* `cis_12` - (Optional) GCP CIS 1.2 Benchmark. Defaults to `false`.
* `k8s` - (Optional) GCP K8S Benchmark. Defaults to `false`.
* `pci` - (Optional) GCP PCI Benchmark. Defaults to `false`.
* `pci_rev2` - (Optional) GCP PCI Benchmark Rev2. Defaults to `false`.
* `soc` - (Optional) GCP SOC 2 Report. Defaults to `false`.
* `soc_rev2` - (Optional) GCP SOC 2 Report Rev2. Defaults to `false`.

### Daily Compliance Reports

`daily_compliance_reports` supports the following arguments:

* `host_security` - (Optional) Host Security. Defaults to `false`.
* `platform` - (Optional) Platform Events. Defaults to `false`.
* `openshift_compliance` - Openshift Compliance (Optional) Defaults to `false`.
* `openshift_compliance_events` - Openshift Compliance Events (Optional) Defaults to `false`.
* `aws_cloudtrail` - (Optional) AWS CloudTrail. Defaults to `false`.
* `aws_compliance` - (Optional) AWS Compliance. Defaults to `false`.
* `azure_activity_log` - (Optional) Azure Activity Log. Defaults to `false`.
* `gcp_audit_trail` - (Optional) GCP Audit Trail. Defaults to `false`.
* `gcp_compliance` - (Optional) GCP Compliance. Defaults to `false`.

## Import

A Lacework Report Rule can be imported using a `GUID`, e.g.

```
$ terraform import lacework_report_rule.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
