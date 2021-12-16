---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_ghcr"
description: |-
  Create and manage Github Container Registry (GHCR) integrations
---

# lacework\_integration\_ghcr

Use this resource to integrate a Github Container Registry (GHCR) with Lacework to assess, identify,
and report vulnerabilities found in the operating system software packages in container images. 
For more information, see the [Integrate Github Container Registry documentation](https://docs.lacework.com/integrate-github-container-registry).

## Example Usage

```hcl
resource "lacework_integration_ghcr" "example" {
  name     = "My Github registry Registry"
  non_os_package_support = true
  username = "my-user"
  password = "a-secret-password"
  ssl      = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The integration name.
* `username` - (Required) The Github username.
* `password` - (Required) The Github personal access token with permission `read:packages`.
* `ssl` - (Optional) Enable or disable SSL communication. Defaults to `true`.
* `registry_notifications` - (Optional) Subscribe to registry notifications. Defaults to `false`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `limit_by_tags` - (Optional) A list of image tags to limit the assessment of images with matching tags. If you specify `limit_by_tags` and `limit_by_label` limits, they function as an `AND`.
* `limit_by_label` - (Optional) A list of key/value labels to limit the assessment of images. If you specify `limit_by_tags` and `limit_by_label` limits, they function as an `AND`.
* `limit_by_repositories` - (Optional) A list of repositories to assess.
* `non_os_package_support` - (Optional) Enable [program language scanning](https://docs.lacework.com/container-image-support#language-libraries-support). Defaults to `true`.

The `limit_by_label` block can be defined multiple times to define multiple label limits, it supports:
* `key` - (Required) The key of the label.
* `value` - (Required) The value of the label.

## Import

A Lacework Github container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_ghcr.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
