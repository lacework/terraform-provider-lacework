---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_gcp_agentless_scanning"
description: |-
Create and manage GCP Organizations Agentless Scanning integration
---

# lacework\_integration\_gcp\_agentless\_scanning

Use this resource to configure a GCP Agentless Scanning integration.

## Example Usage

```hcl
resource "lacework_integration_gcp_agentless_scanning" "account_abc" {
  name         = "Integration name"
  resource_level = "PROJECT"
  resource_id  = "to-scan-gcp-project-id"
  storage_bucket = "gcp storage bucket hosting shared results"
  scanning_project_id = "lacework scanner project id"
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

* `name` - (Required) The GCP Agentless Scanning integration name.
* `resource_level` = (Required) "PROJECT" or "ORGANIZATION"
* `resource_id` - (Required) The organization or project ID.
* `storage_bucket` - (Required) The bucket arn where analysis results are shared with Lacework platform.

* `scan_frequency` - (Optional) How often, in hours, the scan will run - Defaults to 24 hours.
* `query_text` - (Optional) The lql query.
* `filter_list` - (Optional) Comma separated list to include or exclude projects.
* `scan_containers` - (Optional) Whether to includes scanning for containers.
* `scan_host_vulnerabilities` - (Optional) Whether to include scanning for host vulnerabilities.
* `credentials` - (Optional) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `retries` - (Optional) The number of attempts to create the external integration. Defaults to `5`.

### Credentials
These are the credentials of the service account which has read only access to the storage bucket.

`credentials` supports the following arguments:

* `client_id` - (Required) The service account client ID.
* `client_email` - (Required) The service account client email.
* `private_key_id` - (Required) The service account private key ID.
* `private_key` - (Required) The service account private key.

## Import

A Lacework GCP Agentless Scanning integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_gcp_agentless_scanning.account_abc EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
