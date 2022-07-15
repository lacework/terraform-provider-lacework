---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_agentless_scanning"
description: |-
  Create and manage AWS Agentless Scanning integration
---

# lacework\_integration\_aws\_agentless\_scanning

Use this resource to configure an AWS Agentless Scanning integration.

## Example Usage

```hcl
resource "lacework_integration_aws_agentless_scanning" "account_abc" {
  name                      = "account ABC"
  scan_frequency            = 24
  query_text                = var.query_text
  scan_containers           = true
  scan_host_vulnerabilities = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The AWS Config integration name.
* `scan_frequency` - (Required) How often in hours the scan will run.
* `query_text` - (Optional) The lql query.
* `scan_containers` - (Optional) Whether to includes scanning for containers.
* `scan_host_vulnerabilities` - (Optional) Whether to includes scanning for host vulnerabilities.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

## Import

A Lacework AWS Agentless Scanning integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_agentless_scanning.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
