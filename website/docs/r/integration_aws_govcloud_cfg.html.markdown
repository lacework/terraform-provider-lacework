---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_govcloud_cfg"
description: |-
  Create and manage AWS GovCloud Config integrations
---

# lacework\_integration\_aws\_govcloud\_cfg

Use this resource to configure an AWS Config integration for AWS GovCloud to analyze AWS configuration compliance.

To find more information see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360021140214-Initial-Setup-of-AWS-GovCloud-Integration).

## Example Usage

```hcl
resource "lacework_integration_aws_govcloud_cfg" "example" {
	name = "AWS gov cloud config integration example"
	account_id = "553453453"
	credentials {
		access_key_id     = "AWS123abcAccessKeyID"
		secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The AWS GovCloud Config integration name.
* `account_id` - (Required) The AWS account ID.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

`credentials` supports the following arguments:

* `access_key_id` - (Required) The AWS access key ID.
* `secret_access_key` - (Required) The AWS secret key for the specified AWS access key.

## Import

A Lacework AWS Config integration for AWS GovCloud can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_govcloud_cfg.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
