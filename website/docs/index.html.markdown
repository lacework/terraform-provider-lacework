---
layout: "lacework"
page_title: "Provider: Lacework"
description: |-
  The Lacework provider is used to interact with the Lacework cloud security platform.
---

# Lacework Provider

The Lacework provider is used to interact with the Lacework cloud security platform.

## Example Usage

```hcl
# Configure the Lacework Provider
provider "lacework" {
  account     = "${var.lacework_account}"
  api_key     = "${var.lacework_api_key}"
  api_secret  = "${var.lacework_api_secret}"
}

# Connect an AWS account to Lacework for configuration and compliance assessment
resource "lacework_integration_aws_cfg" "account_a" {
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `account` - (Required) This is the Lacework account subdomain of URL (i.e. `<ACCOUNT>`
  .lacework.net). It can also be sourced from the `LW_ACCOUNT` environment variable.

* `api_key` - (Required) This is a Lacework API access key. It can also be sourced
  from the `LW_API_KEY` environment variable.

* `api_key` - (Required) This is a Lacework API access secret. It can also be sourced
  from the `LW_API_SECRET` environment variable.
