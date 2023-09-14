---
subcategory: "Other Resources"
layout: "lacework"
page_title: "Lacework: lacework_external_id"
description: |-
  Generates an External ID using Lacework format.
---

# lacework\_external\_id

This resource generates an External ID (EID) using Lacework format. These IDs are used to create integrations.

The v2 format is:
```
lweid:<csp>:<version>:<lw_tenant_name>:<aws_acct_id>:<random_string_size_10>
```

## Example Usage

```hcl
resource "lacework_external_id" "aws_123456789012" {
  csp        = "aws"
  account_id = "123456789012"
}
```

## Argument Reference

* `csp` - (Required) The Cloud Service Provider. Valid CSP's include: `aws`, `google`, `oci`, and `azure`.
* `account_id` - (Required) The account id from the CSP to be integrated.

## Attribute Reference

The following attributes are exported:

* `v2` - The generated External ID version 2.

