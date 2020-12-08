---
subcategory: "Agents"
layout: "lacework"
page_title: "Lacework: lacework_agent_access_token"
description: |-
  Manage agent access tokens
---

# lacework\_agent\_access\_token

To connect to the Lacework platform, Lacework agents require an agent access token. Use this resource to
mange agent tokens within your Lacework account. 

!> **Warning:** Agent tokens should be treated as secret and not published. A token uniquely identifies
a Lacework customer. If you suspect your token has been publicly exposed or compromised, generate a new
token, update the new token on all machines using the old token. When complete, the old token can safely
be disabled without interrupting Lacework services.

You can use the agent token name to logically separate your deployments, for example, by environment types
(QA, Dev, etc.) or system types (CentOS, RHEL, etc.).

-> **Note:** The Lacework agent runs on most Linux distributions. For more detailed information, see
	[Supported Operating Systems.](https://support.lacework.com/hc/en-us/articles/360005230014-Supported-Operating-Systems).

!> **Warning:** By design, agent tokens cannot be deleted. Running terraform destroy will only disable the token.

## Example Usage

```hcl
resource "lacework_agent_access_token" "k8s" {
  name        = "prod"
  description = "k8s deployment for production env"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The agent access token name.
* `description` - (Optional) The agent access token description.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `token` - The agent access token.

## Import

A Lacework agent access token can be imported using the token itself, e.g.

```
$ terraform import lacework_agent_access_token.k8s YourAgentToken
```
-> **Note:** To list all agent access tokens in your Lacework account, use the
	Lacework CLI command `lacework agent token list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).

