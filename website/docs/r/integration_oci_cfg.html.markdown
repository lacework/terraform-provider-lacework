---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_oci_cfg"
description: |-
  Create and manage OCI Config integrations
---

# lacework\_integration\_oci\_cfg

Use this resource to configure an OCI Configuration integration to analyze OCI compliance.

## Example Usage

```hcl
resource "lacework_integration_oci_cfg" "account_abc" {
  name = "account ABC"
  tenant_id = "ocid1.tenancy.oc1..abcdefghijklmnopqrstuvxyz1234567890"
  tenant_name = "tenant_xyz"
  home_region = "us-sanjose-1"
  user_ocid = "ocid1.user.oc1..abcdefghijklmnopqrstuvxyz1234567890"
  credentials {
    fingerprint = "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f"
    private_key = "-----BEGIN PRIVATE KEY-----\n ... -----END PRIVATE KEY-----\n"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The OCI Configuration integration name.
* `tenant_id` - (Required) The OCID of the tenant to be integrated with Lacework.
* `tenant_name` - (Required) The name of the tenant to be integrated with Lacework.
* `home_region` - (Required) The home region of the tenant to be integrated with Lacework.
* `user_ocid` - (Required) The OCID of the OCI user used used in the integration.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.

### Credentials

`credentials` supports the following arguments:

* `fingerprint` - (Required) The fingerprint of the public key used for authentication.
* `private_key` - (Required) The private key used for authentication in PEM format.

## Import

A Lacework OCI Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_oci_cfg.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework cloud-account list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
