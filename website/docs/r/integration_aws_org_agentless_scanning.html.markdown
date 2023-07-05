---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_org_agentless_scanning"
description: |-
  Create and manage AWS Organizations Agentless Scanning integration
---

# lacework\_integration\_aws\_org\_agentless\_scanning

Use this resource to configure an AWS Organizations Agentless Scanning integration.

## Example Usage

```hcl
resource "lacework_integration_aws_org_agentless_scanning" "account_abc" {
  name                      = "account ABC"
  scan_frequency            = 24
  query_text                = var.query_text
  scan_containers           = true
  scan_host_vulnerabilities = true

  account_id = "0123456789"
  bucket_arn = "arn:aws:s3:::bucket-arn"

  scanning_account   = "0123456789"
  management_account = "0123456789"
  monitored_accounts = ["r-1234"]

	credentials { 
	  role_arn = "arn:aws:iam::0123456789:role/iam-123"
	  external_id = "0123456789"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The AWS Organizations Agentless Scanning integration name.
* `scan_frequency` - (Required) How often, in hours, the scan will run.
* `query_text` - (Optional) The LQL query.
* `scan_containers` - (Optional) Whether to includes scanning for containers.
* `scan_host_vulnerabilities` - (Optional) Whether to includes scanning for host vulnerabilities.
* `scan_multi_volume` - (Optional) Whether to scan secondary volumes (`true`) or only root volumes (`false`). Defaults to `false`
* `scan_stopped_instances` - (Optional) Whether to scan stopped instances (`true`). Defaults to `true`
* `account_id` - (Optional) The AWS account ID.
* `bucket_arn` - (Optional) The bucket ARN.
* `scanning_account` - (Required) The scanning AWS account ID.
* `management_account` - (Optional) The management AWS account ID.
* `monitored_accounts` - (Required) The list of monitroed AWS account IDs or OUs.
* `credentials` - (Optional) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

  `credentials` supports the following arguments:

* `role_arn` - (Optional) The role arn.
* `external_id` - (Optional) The external id.

## Import

A Lacework AWS Organizations Agentless Scanning integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_org_agentless_scanning.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework cloud-accounts list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
