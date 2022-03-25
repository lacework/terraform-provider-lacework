---
subcategory: "Team Members"
layout: "lacework"
page_title: "Lacework: lacework_team_member"
description: |-
  Create and manage Team Members
---

# lacework\_team\_member

Team members can be granted access to multiple Lacework accounts and have different roles for each account. Team members can also be granted organization-level roles.

Lacework supports the following account roles:
* **Administrator Role** - Full access to all functionalities
* **User Role** - View only and no access to API keys

Lacework supports the following organization roles:
* **Org Administrator Role** - The member has admin privileges for organization settings and admin privileges for all accounts within the organization
* **Org User Role** - The member has user privileges for organization settings and user privileges for all accounts within the organization

For more information, see the [Team Members documentation](https://docs.lacework.com/team-members).

## Example Usage: Standalone Account

Team member that has **User Role** access into a Lacework standalone account.

```hcl
resource "lacework_team_member" "harry" {
  first_name = "Harry"
  last_name  = "Potter"
  email      = "harry@hogwarts.io"
}
```
Team member with **Administrator Role** access.

```hcl
resource "lacework_team_member" "hermione" {
  first_name    = "Hermione"
  last_name     = "Granger"
  email         = "hermione@hogwarts.io"
  administrator = true
}
```

## Example Usage: Organizational Account

To manage team members at the organization-level you need to define a Lacework provider with the
`organization` argument set to `true`.

```hcl
provider "lacework" {
  alias        = "org"
  organization = true
}
```

Team member that has **User Role** access for all accounts within the Lacework organization.

```hcl
resource "lacework_team_member" "ron" {
  provider   = lacework.org
  first_name = "Ron"
  last_name  = "Weasley"
  email      = "ron@hogwarts.io"

  organization {
    user = true
  }
}
```

Team member with **Administrator** privileges for all accounts within the organization.

```hcl
resource "lacework_team_member" "albus" {
  provider   = lacework.org
  first_name = "Albus"
  last_name  = "Dumbledore"
  email      = "albus@hogwarts.io"

  organization {
    administrator = true
  }
}
```

Team Member that has access to multiple accounts within a Lacework organization. The member is an
administrator for `SLYTHERIN` account, and a regular user for `HUFFLEPUFF`, `RAVENCLAW`, and `GRYFFINDOR`.

```hcl
resource "lacework_team_member" "severus" {
  provider   = lacework.org
  first_name = "Severus"
  last_name  = "Snape"
  email      = "severus@hogwarts.io"

  organization {
    admin_accounts = ["SLYTHERIN"]
    user_accounts  = ["HUFFLEPUFF", "RAVENCLAW", "GRYFFINDOR"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `first_name` - (Required) The team member first name.
* `last_name` - (Required) The team member last name.
* `email` - (Required) The team member email address, which will also be used as the username.
* `administrator` - (Optional) Set to `true` to make the team member an administrator, otherwise the member will be a regular user. Defaults to `false`.
* `organization` - (Optional) Use this block to manage organization-level team members. See [Organization](#organization) below for details.
* `enabled` - (Optional) The state of the team member. Defaults to `true`.

### Organization

`organization` supports the following arguments:

* `user` - (Optional) Whether the team member is an organization-level user. Defaults to `false`.
* `administrator` - (Optional) Whether the team member is an organization-level administrator. Defaults to `false`.
* `user_accounts` - (Optional) List of accounts the team member is a user.
* `admin_accounts` - (Optional) List of accounts the team member is an administrator.

## Import

There are two ways to import a team member.

### Import Standalone Team Member
A Lacework standalone team member can be imported using a `USER_GUID`, e.g.

```
$ terraform import lacework_team_member.harry HOGWARTS_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```

### Import Organizational Team Member
A Lacework organization-level team member can be imported using the `email`, e.g.

```
$ terraform import lacework_team_member.albus albus@hogwarts.io
```

-> **Note:** To retrieve the `USER_GUID` or `EMAIL` from existing team members in your account,
use the Lacework CLI command `lacework team-member list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
