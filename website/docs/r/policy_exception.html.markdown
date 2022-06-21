---
subcategory: "Policy Exceptions"
layout: "lacework"
page_title: "Lacework: lacework_policy_exception"
description: |-
  Create and manage Lacework Policy Exceptions
---

# lacework\_policy\_exception

Add exceptions to Lacework policies

For more information, see the [Adding Exceptions to a Policy](https://docs.lacework.com/add-exception-to-LQL-policy-lw-console).

## Example Usage

Create a Lacework Policy Exception to exempt specified aws account from policy.

```hcl
resource "lacework_policy_exception" "example" {
  policy_id   = "lacework-global-73"
  description = "Exception for account 123456789"
  constraint {
    field_key   = "accountIds"
    field_values = ["123456789"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) The description of the policy exception.
* `policy_id` - (Required) The id of the policy the exception is associated.
* `constraint` - (Required) Constraint. See [Constraint](#Constraint) below for details.

### Constraint

`constraint` supports the following arguments:

* `field_key` - (Required) The key of the constraint being applied. Example for Aws polices this could be `accountIds`.
* `field_values` - (Required) The values related to the constraint key.

## Import

A Lacework policy can be imported using a `POLICY_ID` and `EXCEPTION_ID`, e.g.

```
$ terraform import lacework_policy.example YourLQLPolicyID ab1234c5-de6f-789g-1234-5hi6789jk1lm
```
