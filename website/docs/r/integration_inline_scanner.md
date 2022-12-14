---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_inline_scanner"
description: |-
  Create and manage Inline Scanner container registry tokens
---

# lacework\_integration\_inline\_scanner

Use this resource to integrate a Inline Scanner with Lacework to assess, identify,
and report vulnerabilities found as part of Inline Scanner integration.

## Example Usage

```hcl
resource "lacework_integration_inline_scanner" "example" {
  name = "My Inline Scanner Example" 
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Container Registry integration name.

* `limit_num_scans` - (Optional) The maximum number of scans per hour that this integration can perform. Defaults to `60`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `integration_tags` - (Optional) Identifier tags as `key:value` pairs.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `server_token` - The inline scanner access token.
* `server_token_uri` - The inline scanner github path.
* `policy_evaluate` - A `bool` value indicating whether a policy is associated to this token.
* `policy_guids` - A `list` policy guids associated to this token.

## Import

A Lacework Inline Scanner container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_inline_scanner.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework container-registry list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).