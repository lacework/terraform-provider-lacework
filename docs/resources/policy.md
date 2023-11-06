---
subcategory: "Policies"
layout: "lacework"
page_title: "Lacework: lacework_policy"
description: |-
  Create and manage Lacework Policies
---

# lacework\_policy

Lacework provides a highly scalable platform for creating, customizing, and managing custom policies
against any datasource that is exposed via the Lacework Query Language (LQL).

For more information, see the [Policy Overview Documentation](https://docs.lacework.net/console/custom-policy-overview).

## Example Usage

Create a Lacework Policy to check for a change of password from an RDS cluster.

```hcl
resource "lacework_query" "AWS_CTA_AuroraPasswordChange" {
  query_id = "TF_AWS_CTA_AuroraPasswordChange"
  query    = <<EOT
  {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 'rds.amazonaws.com'
          and EVENT_NAME = 'ModifyDBCluster'
          and value_exists(EVENT:requestParameters.masterUserPassword)
          and EVENT:requestParameters.applyImmediately = true
          and ERROR_CODE is null
      }
      return distinct {
          INSERT_ID,
          INSERT_TIME,
          EVENT_TIME,
          EVENT
      }
  }
EOT
}

resource "lacework_policy" "example" {
  title       = "Aurora Password Change"
  description = "Password for an Aurora RDS cluster was changed"
  remediation = "Check that the password change was expected and ensure only specified users can modify the RDS cluster"
  query_id    = lacework_query.AWS_CTA_AuroraPasswordChange.id
  severity    = "High"
  type        = "Violation"
  evaluation  = "Hourly"
  tags        = ["cloud_AWS", "custom"]
  enabled     = false

  alerting {
    enabled = false
    profile = "LW_CloudTrail_Alerts.CloudTrailDefaultAlert_AwsResource"
  }
}
```

-> **Note:** Lacework automatically generates a policy id when you create a policy, which is the recommended workflow.
Optionally, you can define your own policy id using the `policy_id_suffix`, this suffix must be all lowercase letters,
optionally followed by `-` and numbers, for example, `abcd-1234`. When you define your own policy id, Lacework prepends
the account name. The final policy id would then be `lwaccountname-abcd-1234`.

## Argument Reference

The following arguments are supported:

* `title` - (Required) The policy title.
* `description` - (Required) The description of the policy.
* `query_id` - (Required) The query id.
* `severity` - (Required) The list of the severities. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.
* `type` - (Required) The policy type must be `Violation`.
* `evaluation` - (Optional) The evaluation frequency at which the policy will be evaluated. Valid values are
  `Hourly` or `Daily`. Defaults to `Hourly`.
* `remediation` - (Required) The remediation message to display.
* `limit` - (Optional) Set the maximum number of records returned by the policy.
   Maximum value is `5000`. Defaults to `1000`
* `enabled` - (Optional) Whether the policy is enabled or disabled. Defaults to `true`.
* `policy_id_suffix` - (Optional) The string appended to the end of the policy id.
* `tags` - (Optional) A list of policy tags.
* `alerting` - (Optional) Alerting. See [Alerting](#alerting) below for details.

### Alerting

`alerting` supports the following arguments:

* `profile` - (Required) The alerting profile.
* `enabled` - (Optional) Whether the alerting profile is enabled or disabled. Defaults to `true`.

## Import

A Lacework policy can be imported using a `POLICY_ID`, e.g.

```
$ terraform import lacework_policy.example YourLQLPolicyID
```

-> **Note:** To retrieve the `POLICY_ID` from existing policies in your account, use the
Lacework CLI command `lacework policy list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
