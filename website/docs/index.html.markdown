---
layout: "lacework"
page_title: "Provider: Lacework"
description: |-
  The Lacework provider is used to interact with the Lacework cloud security platform.
---

# Lacework Provider

The Lacework provider is used to interact with the Lacework cloud security platform.
The provider needs to be configured with the proper credentials before it can be used.

Use the left navigation panel to read about the available resources.

## Example Usage

```hcl
# Configure the Lacework Provider
provider "lacework" {
  profile = "my-profile"
}

# Connect an AWS account to Lacework for configuration and compliance assessment
resource "lacework_integration_aws_cfg" "account_abc" {
  # ...
}

# Configure Lacework to forward alerts to a Slack channel
resource "lacework_alert_channel_slack" "critical" {
  # ...
}
```

## Authentication
The Lacework provider can be configured with the proper credentials via the following supported methods:

* Static credentials
* Environment variables
* Configuration file

### Static credentials
!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended. Secrets could be leaked by committing this file to a public version
control system.

Static credentials can be provided by adding the `account`, `api_key`, and `api_secret` in-line in the
Lacework provider block:

```hcl
provider "lacework" {
  account    = "my-account"
  api_key    = "my-api-key"
  api_secret = "my-api-secret"
}
```

### Environment Variables
You can provide your credentials via the `LW_ACCOUNT`, `LW_API_KEY`, and `LW_API_SECRET` environment
variables, they represent your Lacework account subdomain of URL, Lacework API access key, and Lacework
API access secret, respectively.

-> **Note:** Setting your Lacework credentials using these environment variables will override the use of `LW_PROFILE`.

```hcl
provider "lacework" {}
```

Terminal:

```
$ export LW_ACCOUNT="my-account"
$ export LW_API_KEY="my-api-key"
$ export LW_API_SECRET="my-api-secret"
$ terraform plan
```

### Configuration file
It is possible to use credentials from the Lacework configuration file. The default location on Linux and OS X
is `$HOME/.lacework.toml`, and for Windows users is `"%USERPROFILE%\.lacework.toml"`. This configuration file
can be easily managed using the [Lacework CLI](https://github.com/lacework/go-sdk/wiki/CLI-Documentation). This
method also supports a `profile` configuration and matching `LW_PROFILE` environment variable.

```hcl
provider "lacework" {
  profile = "custom-profile"
}
```

## Organizational Accounts

An organization can contain multiple accounts so you can manage components such as alerts, resource groups,
team members, and audit logs at a more granular level inside an organization. A team member may have access
to multiple accounts and can easily switch between them.

-> **Note:** To enroll your account in an organization follow the [Organization Enrollment Process](https://support.lacework.com/hc/en-us/articles/360041727394-Organization-Overview) documentation.

Use the argument `subaccount` to switch to a different account inside your organizational account.

For example, having a `default` profile that has access to your primary account named `my-company`:
```toml
[default]
  account = "my-company"
  api_key = "my-api-key"
  api_secret = "my-api-secret"
  version = 2
```

To access your sub-account named `business-unit` you would specify the argument `subaccount`.
```hcl
provider "lacework" {
  alias = "primary"
}

provider "lacework" {
  alias      = "business-unit"
  subaccount = "business-unit"
}
```

From there, you can pass the [`alias` meta-argument](https://www.terraform.io/docs/language/providers/configuration.html#alias-multiple-provider-configurations) to any resource to switch between accounts:
```hcl
resource "lacework_alert_channel_slack" "primary_critical" {
  alias = lacework.primary
  # ...
}
resource "lacework_alert_channel_slack" "business_unit_critical" {
  alias = lacework.business-unit
  # ...
}
```

!> **Warning:** To manage multiple accounts, your user should have the Organization Administrator Role.

### Organization Level Access

Organization administrators can access organization level data sets by setting the `organization` argument to `true`.
```hcl
provider "lacework" {
  organization = true
}
```

!> **Warning:** When accessing organization level data sets, the `subaccount` argument is ignored.

Using this type of configuration is intended for managing resources such as alerts, resource groups, team members,
cloud accounts, and more, at the organization level.

### Migrating existing resources to the Organization level

When attempting to migrate an existing resource from one of your Lacework accounts to the organization level,
you need to delete the resource, update the Lacework provider to access the organization level data set, and
run `terraform apply` to create a new resource at the organization level.

## Argument Reference

The following arguments are supported in the `provider` block:

* `profile` - (Optional) This is the Lacework profile name to use, profiles are configured
  at `$HOME/.lacework.toml` via the [Lacework CLI](https://github.com/lacework/go-sdk/wiki/CLI-Documentation).
  It can also be sourced from the `LW_PROFILE` environment variable.

* `account` - (Optional) This is the Lacework account subdomain of URL (i.e. `<ACCOUNT>`
  .lacework.net). It must be provided, but it can also be sourced from the `LW_ACCOUNT`
  environment variable, or via the configuration file if `profile` is specified.

* `api_key` - (Optional) This is a Lacework API access key. It must be provided, but it can
  also be sourced from the `LW_API_KEY` environment variable, or via the configuration file
  if `profile` is specified.

* `api_secret` - (Optional) This is a Lacework API access secret. It must be provided, but it
  can also be sourced from the `LW_API_SECRET` environment variable, or via the configuration
  file if `profile` is specified.

* `subaccount` - (Optional) The sub-account name inside your organization (for organization
  administrators only). It can also be sourced from the `LW_SUBACCOUNT` environment variable,
  or via the configuration file if `profile` is specified.

* `organization` - (Optional) Set this argument to `true` to access organization level data
  sets (for organization administrators only). It can also be sourced from the `LW_ORGANIZATION`
  environment variable.

-> **Note:** To generate a set of API access keys follow [this documentation](https://support.lacework.com/hc/en-us/articles/360011403853-Generate-API-Access-Keys-and-Tokens).
