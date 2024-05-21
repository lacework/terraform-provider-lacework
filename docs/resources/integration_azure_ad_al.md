---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_azure_ad_al"
description: |-
  Create and manage Azure Active Directory Activity Log integrations
---

# lacework\_integration\_azure\_ad\_al

!> **Warning:** This integration is not yet generally available. Please contact your Lacework account team to request access to the Azure AD feature preview.

Use this resource to configure an Azure Active Directory Activity Log integration to analyze audit logs
for monitoring cloud account security.

## Example Usage

```hcl
resource "lacework_integration_azure_ad_al" "account_abc" {
  name                = "account ABC"
  tenant_id           = "abbc1234-abc1-123a-1234-abcd1234abcd"
  event_hub_namespace = "your-eventhub-ns.servicebus.windows.net"
  event_hub_name      = "your-event-hub-name"
  credentials {
    client_id     = "1234abcd-abcd-1234-ab12-abcd1234abcd"
    client_secret = "ABCD1234abcd1234abdc1234ABCD1234abcdefxxx="
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Azure Active Directory Activity Log integration name.
* `tenant_id` - (Required) The directory tenant ID.
* `event_hub_namespace` - (Required) The EventHub Namespace.
* `event_hub_name` - (Required) The EventHub Name.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The application client ID.
* `client_secret` - (Required) The client secret.

## Import

A Lacework Azure Active Directory Activity Log integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_azure_ad_al.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework cloud-account list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
