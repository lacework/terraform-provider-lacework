---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_cloud_account_aws_eks_audit_log"
description: |-
  Create and manage AWS EKS Audit Log integrations
---

# lacework\_cloud\_account\_aws\_eks\_audit\_log

Use this resource to configure an AWS EKS Audit Log integration to analyze EKS audit logs.

## Example Usage

```hcl
resource "lacework_cloud_account_aws_eks_audit_log" "account_abc" {
  name      = "account ABC"
  sns_arn   = "arn:aws:sns:us-west-2:123456789:foo-lacework-eks:00777777-ab77-1234-a123-a12ab1d12c1d"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The AWS CloudTrail integration name.
* `sns_arn` - (Required) The SNS ARN.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

`credentials` supports the following arguments:

* `role_arn`: (Required) The ARN of the IAM role.
* `external_id`: (Required) The external ID for the IAM role.

## Import

A Lacework AWS Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_ct.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
