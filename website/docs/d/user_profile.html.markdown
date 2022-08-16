---
subcategory: "User Profile"
layout: "lacework"
page_title: "Lacework: lacework_user_profile"
description: |-
  Fetch the current user's Lacework profile.
---

# lacework\_user\_profile

Use this data source to retrieve the User Profile of the current Lacework user.

## Example Usage

```hcl
data "lacework_user_profile" "current" { }
```

## Attribute Reference

The following attributes are exported:

* `username` - The username of the current user.
* `org_account` - A boolean representing whether the user has an organization account.
* `org_admin` - A boolean representing whether the user is an organization admin.
* `org_user` - A boolean representing whether the user is an organization user.
* `url` - A string representing the login URL for the Lacework account.
* `accounts` - An array of accounts in the Lacework tenant.
    * `account_name` - A string representing the account name.
    * `admin` - A boolean representing whether the user is an account admin.
    * `cust_guid` - A string representing the Customer GUID for the account.
    * `user_enabled` - A boolean representing whether the user is enabled.
    * `user_guid` - A string representing the User GUID in the account.
