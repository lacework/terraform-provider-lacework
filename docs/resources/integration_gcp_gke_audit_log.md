---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_gcp_gke_audit_log"
description: |-
  Create and manage GCP GKE Audit Log integrations
---

# lacework\_integration\_gcp\_gke\_audit\_log

Use this resource to configure an [GCP GKE Audit Log integration](https://docs.lacework.com/category/gke-audit-log-integrations) to analyze GKE audit logs.

## Example Usage

```hcl
resource "lacework_integration_gcp_gke_audit_log" "account_abc" {
	name             = "account ABC"
	project_id       = "ABC-project-id"
	subscription     = "projects/ABC-project-id/subscriptions/example-subscription"
	integration_type = "PROJECT"
	credentials {
		client_id      = "123456789012345678900"
		client_email   = "email@abc-project-name.iam.gserviceaccount.com"
		private_key_id = "1234abcd1234abcd1234abcd1234abcd1234abcd"
		private_key    = "-----BEGIN PRIVATE KEY-----\n ... -----END PRIVATE KEY-----\n"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The GCP Audit Trail integration name.
* `organization_id` - (Optional) The organization ID. Required if `integration_type` is set to `ORGANIZATION`.
* `project_id` - (Required) The project ID.
* `subscription` - (Required) The PubSub Subscription.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `integration_type` - (Optional) The integration type. Must be one of `PROJECT` or `ORGANIZATION`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The service account client ID.
* `client_email` - (Required) The service account client email.
* `private_key_id` - (Required) The service account private key ID.
* `private_key` - (Required) The service account private key.

## Import

A Lacework GCP GKE Audit Log integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_gcp_gke_audit_log.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework cloud-account list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
