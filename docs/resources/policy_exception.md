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


Create a Lacework Policy Exception to exempt specified `resourceTags` from policy.

```hcl
resource "lacework_policy_exception" "example" {
  policy_id   = "lacework-global-73"
  description = "Exception for resource tag example1 and example2"

  constraint {
    field_key = "resourceTags"
    field_values_map {
      key   = "example_tag1"
      value = ["example_value", "example_value1"]
    }
    field_values_map {
      key   = "example_tag2"
      value = ["example_value", "example_value1"]
    }
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
* `field_values` - (Optional) The values related to the constraint key.
* `field_value_map` - (Optional, **Deprecated**) FieldValueMap. See[FieldValueMap](#FieldValueMap) below for details.
* `field_values_map` - (Optional) FieldValueMap. See[FieldValuesMap](#FieldValuesMap) below for details.

### FieldValueMap

`field_value_map` allows defining constraint values for the `resourceTags` field key. Where `field_value_map` key
property is the name a given resource tag, and `value` includes any value that should match this exception. 
**Deprecated** Use `field_values_map` instead.

`field_values_map` allows defining constraint values for the `resourceTags` field key. Where `field_values_map` key
property is the name a given resource tag, and `value` includes a list of values that should match this exception.

## Import

A Lacework policy can be imported using a `POLICY_ID` and `EXCEPTION_ID`, e.g.

```
$ terraform import lacework_policy_exception.example YourLQLPolicyID YourExceptionID
```
