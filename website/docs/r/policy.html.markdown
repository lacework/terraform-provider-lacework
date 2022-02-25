---
subcategory: "Policy"
layout: "lacework"
page_title: "Lacework: lacework_policy"
description: |-
  Create and manage Policies
---

# lacework\_policy

Lacework provides a highly scalable platform for creating, customizing, and managing custom policies
against any datasource that is exposed via the Lacework Query Language (LQL).

For more information, see the [Policy Overview Documentation](https://docs.lacework.com/custom-policies-overview).

## Example Usage

```hcl
   resource "lacework_query" "example" {
  query_id       = "Lql_Terraform_Query"
  evaluator_id   = "Cloudtrail"
  query          = <<EOT
    Lql_Terraform_Query {
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
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
  title        = "My Policy"
  query_id     = lacework_query.example.id
  severity     = "high"
  type         = "Violation"
  description  = "Policy Created via Terraform"
  remediation  = "Please investigate"
  evaluation   = "Hourly"
  evaluator_id = "Cloudtrail"
  enabled      = true

  alerting {
    enabled = false
    profile = "LW_CloudTrail_Alerts"
  }
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The policy title.
* `query_id` - (Required) The query id.
* `severity` - (Required) The list of the severities. Valid severities include:
  `Critical`, `High`, `Medium`, `Low` and `Info`.
* `type` - (Required) The policy type must be either `Violation` or `Summary`.
* `description` - (Required) The description of the policy.
* `evaluation` - (Required) Set the evaluation frequency `Hourly` or `Daily`.
* `evaluator_id` - (Optional) The evaluator id. `Cloudtrail` must be set for all CloudTrail queries.
* `remediation` - (Optional) The remediation message to display.
* `limit` - (Optional) Set the maximum number of records returned by the policy. Maximum value is `1000`.
* `enabled` - (Optional) Whether the policy is enabled or disabled.
* `policy_id_suffix` - (Optional) The string appended to the end of the policy id.
* `alerting` - (Optional) Alerting. See [Alerting](#alerting) below for details.

### Alerting

`alerting` supports the following arguments:

* `profile` - (Required) The alerting profile.
* `enabled` - (Optional) Whether the alerting profile is enabled or disabled.

## Import

A Lacework policy can be imported using a `POLICY_ID`, e.g.

```
$ terraform import lacework_policy.example MyLQLPolicyID
```

-> **Note:** To retreive the `POLICY_ID` from existing policies in your account, use the
Lacework CLI command `lacework policy list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).