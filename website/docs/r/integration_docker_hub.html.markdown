---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_docker_hub"
description: |-
  Create and manage Docker Hub container registry integrations
---

# lacework\_integration\_docker\_hub

Use this resource to integrate a Docker Hub container registry with Lacework to assess, identify,
and report vulnerabilities found in the operating system software packages in a Docker container
image.

## Example Usage

```hcl
resource "lacework_integration_docker_hub" "example" {
  name = "My Docker Hub Registry Example"
  username = "my-user"
  password = "a-secret-password"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Container Registry integration name.
* `username` - (Required) The Docker user that has at least read-only permissions to the Docker Hub container repositories.
* `password` - (Required) The password for the specified Docker Hub user.
* `limit_by_tag` - (Optional, **Deprecated**) An image tag to limit the assessment of images with matching tag. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. This attribute will be replaced by a new attribute `limit_by_tags` in version 1.0 of the Lacework provider.
* `limit_by_label` - (Optional, **Deprecated**) An image label to limit the assessment of images with matching label. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. This attribute will be replaced by a new attribute `limit_by_labels` in version 1.0 of the Lacework provider.
* `limit_by_repos` - (Optional, **Deprecated**) A comma-separated list of repositories to assess. (without spaces recommended) This attribute will be replaced by a new attribute `limit_by_repositories` in version 1.0 of the Lacework provider.
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `limit_by_tags` - (Optional) A list of image tags to limit the assessment of images with matching tags. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_labels` - (Optional) A key based map of labels to limit the assessment of images with matching `key:value` labels. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_repositories` - (Optional) A list of repositories to assess.

## Import

A Lacework Docker Hub container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_docker_hub.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
