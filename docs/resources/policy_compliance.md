---
subcategory: "Policies"
layout: "lacework"
page_title: "Lacework: lacework_policy_compliance"
description: |-
  Create and manage Lacework Compliance Policies
---

# lacework\_policy\_compliance

Lacework provides a highly scalable platform for creating, customizing, and managing custom policies
against any datasource that is exposed via the Lacework Query Language (LQL).

For more information, see the [Policy Overview Documentation](https://docs.lacework.net/console/custom-policy-overview).

## Example Usage

Create a Lacework Compliance Policy to check for unenabled CloudTrail log file validation.

```hcl
resource "lacework_query" "AWS_Config_CloudTrailLogFileValidationNotEnabled" {
  query_id = "LW_Global_AWS_Config_CloudTrailLogFileValidationNotEnabled"
  query    = <<EOT
  {
    source {
      LW_CFG_AWS_CLOUDTRAIL
    }
    filter {
      RESOURCE_CONFIG:LogFileValidationEnabled = 'false'
    }
    return distinct {
      ACCOUNT_ALIAS,
      ACCOUNT_ID,
      ARN as RESOURCE_KEY,
      RESOURCE_REGION,
      RESOURCE_TYPE,
      SERVICE,
      'CloudTrailLogFileValidationNotEnabled' as COMPLIANCE_FAILURE_REASON
    }
  }
EOT
}

resource "lacework_policy_compliance" "example" {
  query_id = lacework_query.AWS_Config_CloudTrailLogFileValidationNotEnabled.id
  title = "Ensure CloudTrail log file validation is enabled"
  enabled = false
  severity = "High"
  description = "CloudTrail log file validation creates a digitally signed digest\nfile containing a hash of each log that CloudTrail writes to S3. These digest\nfiles can be used to determine whether a log file was changed, deleted, or unchanged\nafter CloudTrail delivered the log. It is recommended that file validation be\nenabled on all CloudTrails."
  remediation = "Perform the following to enable log file validation on a given trail:\nFrom Console:\n1. Sign in to the AWS Management Console and open the IAM console at (https://console.aws.amazon.com/cloudtrail)\n2. Click on Trails on the left navigation pane\n3. Click on target trail\n4. Within the S3 section click on the edit icon (pencil)\n5. Click Advanced\n6. Click on the Yes radio button in section Enable log file validation\n7. Click Save\nFrom Command Line:\naws cloudtrail update-trail --name <trail_name> --enable-log-file-validation\nNote that periodic validation of logs using these digests can be performed by running the following command:\naws cloudtrail validate-logs --trail-arn <trail_arn> --start-time <start_time> --end-time <end_time>"
  tags = ["security:compliance"]
  alerting_enabled = false
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
* `remediation` - (Required) The remediation message to display.
* `enabled` - (Optional) Whether the policy is enabled or disabled. Defaults to `true`.
* `policy_id_suffix` - (Optional) The string appended to the end of the policy id.
* `tags` - (Optional) A list of policy tags.
* `alerting_enabled` - (Optional, **Deprecated**) Whether the alerting profile is enabled or disabled. Defaults to `true`.

## Import

A Lacework compliance policy can be imported using a `POLICY_ID`, e.g.

```
$ terraform import lacework_policy_compliance.example YourLQLPolicyID
```

-> **Note:** To retrieve the `POLICY_ID` from existing policies in your account, use the
Lacework CLI command `lacework policy list`. To install this tool follow
[this documentation](https://docs.lacework.net/cli/).
