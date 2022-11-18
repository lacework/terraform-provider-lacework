---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_proxy_scanner"
description: |-
  Create and manage Proxy Scanner container registry tokens
---

# lacework\_integration\_inline\_scanner

Use this resource to integrate a Proxy Scanner with Lacework to assess, identify,
and report vulnerabilities found as part of Proxy Scanner integration.

## Example Usage

```hcl
resource "lacework_integration_proxy_scanner" "example" {
  name = "My Proxy Scanner Example"
}
```

## Attributes Reference

The following arguments are supported:

* `name` - (Required) The Container Registry integration name.
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `limit_by_tags` - (Optional) A list of image tags to limit the assessment of images with matching tags. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_labels` - (Optional) A key based map of labels to limit the assessment of images with matching `key:value` labels. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_repositories` - (Optional) A list of repositories to assess.
* `policy_evaluate` - A `bool` value indicating whether a policy is associated to this token.
* `policy_guids` - A `list` policy guids associated to this token.

## Argument Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `server_token` - The proxy scanner access token.
* `server_token_uri` - The proxy scanner github path.
## Import

A Lacework Proxy Scanner container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_inline_scanner.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework container-registry list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
