---
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_ct"
description: |-
  Create an manage AWS CloudTrail integrations
---

# lacework\_integration\_aws\_ct

Use this resource to configure an AWS CloudTrail integration to analyze CloudTrail
activity for monitoring cloud account security.

## Example Usage

```hcl
resource "lacework_integration_aws_ct" "account_abc" {
  name      = "account ABC"
  queue_url = "https://sqs.us-west-2.amazonaws.com/123456789012/my_queue"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Lacework AWS CloudTrail integration name.
* `queue_url` - (Required) The SQS Queue URL.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `role_arn`: (Required) The ARN of the IAM role.
* `external_id`: (Required) The external ID for the IAM role.

## Import

A Lacework AWS Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_ct.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/blob/master/cli/README.md).
