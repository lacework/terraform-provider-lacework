---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_gcp_cfg"
description: |-
  Create and manage Google Cloud Config integrations
---

# lacework\_integration\_gcp\_cfg

Use this resource to configure a GCP Config integration to analyze configuration compliance.

## Example Usage

```hcl
resource "lacework_integration_gcp_cfg" "account_abc" {
  name        = "account ABC"
  resource_id = "ABC-project-id"
  credentials {
    client_id      = "123456789012345678900"
    client_email   = "email@abc-project-name.iam.gserviceaccount.com"
    private_key_id = "1234abcd1234abcd1234abcd1234abcd1234abcd"
    private_key    = "-----BEGIN PRIVATE KEY-----\n ... -----END PRIVATE KEY-----\n"
  }
}
```

## GCP Config Module Usage

Lacework maintains a Terraform module that can be used to create and manage the necessary
resources required for both, the cloud provider platform as well as the Lacework platform.

Here is a basic usage of this module:

```hcl
module "config" {
  source  = "lacework/config/gcp"
  version = "~> 0.1.1"
}
```

To see the list of inputs, outputs and dependencies, visit the [Terraform registry page of this module](https://registry.terraform.io/modules/lacework/config/gcp/latest).

## Argument Reference

The following arguments are supported:

* `name` - (Required) The GCP Config integration name.
* `resource_id` - (Required) The organization or project ID.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `resource_level` - (Optional) The integration level. Must be one of `PROJECT` or `ORGANIZATION`. Defaults to `PROJECT`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The service account client ID.
* `client_email` - (Required) The service account client email.
* `private_key_id` - (Required) The service account private key ID.
* `private_key` - (Required) The service account private key.

## Import

A Lacework GCP Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_gcp_cfg.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
