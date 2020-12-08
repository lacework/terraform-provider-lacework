---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_gcr"
description: |-
  Create and manage GCR integrations
---

# lacework\_integration\_gcr

Use this resource to integrate a Google Container Registry (GCR) with Lacework to assess, identify,
and report vulnerabilities found in the operating system software packages in a Docker container
image.

## Example Usage

```hcl
resource "lacework_integration_gcr" "example" {
  name            = "GRC Example"
  registry_domain = "gcr.io"
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

* `name` - (Required) The GCR integration name.
* `registry_domain` - (Required) The GCR domain, which specifies the location where you store the images. Supported domains are `gcr.io`, `us.gcr.io`, `eu.gcr.io`, or `asia.gcr.io`.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `limit_by_tag` - (Optional) An image tag to limit the assessment of images with matching tag. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`.
* `limit_by_label` - (Optional) An image label to limit the assessment of images with matching label. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`.
* `limit_by_repos` - (Optional) A comma-separated list of repositories to assess. (without spaces recommended)
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The service account client ID.
* `client_email` - (Required) The service account client email.
* `private_key_id` - (Required) The service account private key ID.
* `private_key` - (Required) The service account private key.

~> **Note:** The service account used for this integration requires the `storage.objectViewer` role for access to the Google project that contains the Google Container Registry (GCR). The role can be granted at the project level or the bucket level. If granting the role at the bucket level, you must grant the role to the default bucket called `artifacts.[YourProjectID].appspot.com`. In addition, the client must have access to the Google Container Registry API and billing must be enabled.

## Import

A Lacework GCR integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_gcr.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).

