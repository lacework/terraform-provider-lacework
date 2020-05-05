---
layout: "lacework"
page_title: "Lacework: lacework_api_token"
description: |-
  Generate API access token.
---

# lacework\_api\_token

Use this data source to generate API access tokens that can be used to interact with the
external Lacework API.

## Example Usage

```hcl
data "lacework_api_token" "development" { }
```

## Attribute Reference

The following attributes are exported:

* `token` - An API access token.
