---
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_cfg"
description: |-
  Create an manage AWS Config integrations
---

# lacework\_integration\_aws\_cfg

Use this resource to configure an AWS Config integration to analyze AWS configuration compliance.

## Example Usage

```hcl
resource "lacework_integration_aws_cfg" "account_abc" {
  name = "account ABC"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Lacework AWS Config integration name.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `role_arn`: (Required) The ARN of the IAM role.
* `external_id`: (Required) The external ID for the IAM role.

## Import

A Lacework AWS Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_cfg.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/blob/master/cli/README.md).
