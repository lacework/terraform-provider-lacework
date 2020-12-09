---
subcategory: "Agents"
layout: "lacework"
page_title: "Lacework: lacework_agent_access_token"
description: |-
  Lookup agent access token.
---

# lacework\_agent\_access\_token

Retrieve Lacework agent access tokens.

-> **Note:** To list all agent access tokens in your Lacework account, use the
	Lacework CLI command `lacework agent token list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).

## Example Usage

```hcl
data "lacework_agent_access_token" "k8s" {
  name = "k8s-deployments"
}
```

## Argument Reference

* `name` - (Required) The agent access token name.

## Attribute Reference

The following attributes are exported:

* `token` - The agent access token.
