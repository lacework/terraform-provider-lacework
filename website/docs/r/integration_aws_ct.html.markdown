---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_ct"
description: |-
  Create and manage AWS CloudTrail integrations
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

If your Lacework accounts are enrolled in a Lacework organization, you can configure a
consolidated AWS CloudTrail integration that maps CloudTrail activity from your AWS accounts
to selected Lacework accounts within your organization.

The following snippet adds an AWS CloudTrail integration at the Lacework organization level with
the following distribution from AWS accounts to Lacework accounts:

* AWS accounts `234556677` and `774564564` will appear in the Lacework account `lw_acocunt_2` 
* AWS accounts `553453453` and `934534535` will appear in the Lacework account `lw_acocunt_3` 
* All other AWS accounts that are not mapped will appear in the Lacework account `lw_acocunt_1`

```hcl
resource "lacework_integration_aws_ct" "consolidated" {
  name      = "Consolidated CloudTrail"
  queue_url = "https://sqs.us-west-2.amazonaws.com/123456789012/my_queue"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }

  org_account_mappings {
    default_lacework_account = "lw_account_1"

    mapping {
      lacework_account = "lw_account_2"
      aws_accounts     = ["234556677", "774564564"]
    }

    mapping {
      lacework_account = "lw_account_3"
      aws_accounts     = ["553453453", "934534535"]
    }
  }
}
```

-> **Note:** The mapping that you configure for an organization integration is in addition
	to what is already configured for the CloudTrail account integration. It doesn't
	override the existing account integration.

For more information see [Setup of Organization AWS CloudTrail Integration](https://support.lacework.com/hc/en-us/articles/360055993554-Setup-of-Organization-AWS-CloudTrail-Integration)


## Argument Reference

The following arguments are supported:

* `name` - (Required) The AWS CloudTrail integration name.
* `queue_url` - (Required) The SQS Queue URL.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `org_account_mappings` - (Optional) Mapping of AWS accounts to Lacework accounts within a Lacework organization. See [Account Mappings](#organization-account-mappings) below for details.

### Credentials

`credentials` supports the following arguments:

* `role_arn`: (Required) The ARN of the IAM role.
* `external_id`: (Required) The external ID for the IAM role.

### Organization Account Mappings

`org_account_mappings` supports the following arguments:

* `default_lacework_account`: (Required) The default Lacework account name where any non-mapped AWS account will appear.
* `mapping`: (Required) A map of AWS accounts to Lacework account. This can be specified multiple times to map multiple Lacework accounts. See [Mapping](#mapping) below for details.

#### Mapping

The `mapping` block supports:

* `lacework_account`: (Required) The Lacework account name where the CloudTrail activity from the selected AWS accounts will appear.
* `aws_accounts`: (Required) The list of AWS account IDs to map.

## Import

A Lacework AWS Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_aws_ct.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
