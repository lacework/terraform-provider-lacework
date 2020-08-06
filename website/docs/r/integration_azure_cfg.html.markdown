---
layout: "lacework"
page_title: "Lacework: lacework_integration_azure_cfg"
description: |-
  Create and manage Azure Cloud Config integrations
---

# lacework\_integration\_azure\_cfg

Use this resource to configure an Azure Config integration to analyze configuration compliance.

## Example Usage

```hcl
resource "lacework_integration_azure_cfg" "account_abc" {
  name = "account ABC"
  tenant_id = "abbc1234-abc1-123a-1234-abcd1234abcd"
  credentials {
    client_id     = "1234abcd-abcd-1234-ab12-abcd1234abcd"
    client_secret = "ABCD1234abcd1234abdc1234ABCD1234abcdefxxx="
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Lacework Azure Config integration name.
* `tenant_id` - (Required) The directory tenant ID.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The application client ID.
* `client_secret` - (Required) The client secret.

## Import

A Lacework Azure Config integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_azure_cfg.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/blob/master/cli/README.md).

