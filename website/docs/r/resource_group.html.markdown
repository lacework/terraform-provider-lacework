---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: (beta) lacework_resource_group"
description: |-
  Create and manage Resource Groups V2 (Beta)
---

# (beta) lacework\_resource\_group

Use this resource to create a V2 Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://lwdocs-rg2.netlify.app/console/resource-groups/).


## Example Usage

```hcl
resource "lacework_resource_group" "example" {
  name        = "My Resource Group"
  type        = "AWS"
  description = "This groups a subset of AWS resources"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Region"
      operation = "EQUALS"
      value     = ["us-east-1"]
    }

    filter {
      filter_name = "filter2"
      field     = "Region"
      operation = "EQUALS"
      value     = ["us-west-2"]
    }

    group {
      operator = "AND"

      filter {
        filter_name = "filter5"
        field     = "Region"
        operation = "EQUALS"
        value     = ["us-central-1"]
      }

      group {
        operator = "OR"
        filter {
          filter_name = "filter3"
          field     = "Account"
          operation = "EQUALS"
          value     = ["987654321"]
        }
        filter {
          filter_name = "filter4"
          field     = "Account"
          operation = "EQUALS"
          value     = ["123456789"]
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group name.
* `group` - (Required) The representation of the expression that a resource must match to be 
  part of the resource group. Groups can be nested up to 3 levels deep and can be combined by 
  individual filters. See the [api-docs](https://lwdocs-rg2.netlify.app/api/api-resource-group/#filterable-fields) for the supported fields.
  Each `group` must have at least one of `group` or `filter` defined.
* `type` - (Required) The type of resource group being created. e.g. AWS/GCP/AZURE
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Resource Group can be imported using a `RESOURCE_GROUP_GUID`, e.g.

```
$ terraform import lacework_resource_group.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `RESOURCE_GROUP_GUID` from existing resource groups in your account, 
use the Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
