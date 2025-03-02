---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: lacework_resource_group"
description: |-
  Create and manage Resource Groups in Lacework
---

# lacework\_resource\_group

Use this resource to create a Resource Group in order to categorize Lacework-identifiable assets.
For more information, see the [Resource Groups documentation](https://docs.fortinet.com/document/lacework-forticnapp/latest/api-reference/690087/using-the-resource-groups-api).

## Converting Original to Newer Resource Groups

Please refer to this [documentation](https://docs.fortinet.com/document/lacework-forticnapp/latest/api-reference/375795/convert-original-to-newer-resource-groups-in-terraform) to understand how to
convert the original resource groups to the newer resource groups.

## Example Usage

The following Terraform code defines a Lacework resource group that includes all AWS resources that are located in us-east-1 or us-west-2 or those in us-central-1 with an account ID of either 987654321 or 123456789.

```hcl
resource "lacework_resource_group" "example" {
  name        = "My Resource Group"
  type        = "AWS"
  description = "This groups a subset of AWS resources"
  group {
    operator = "OR"
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
  individual filters. See the [api-docs](https://docs.fortinet.com/document/lacework-forticnapp/latest/api-reference/690087/using-the-resource-groups-api#filterable-fields) for the supported fields.
  Each `group` must have at least one of `group` or `filter` defined.
* `type` - (Required) The type of resource group being created, AWS, GCP, or AZURE
* `description` - (Optional) The description of the resource group.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

You can import a Lacework resource group by `RESOURCE_GROUP_GUID`, for example:

```
$ terraform import lacework_resource_group.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```

-> **Note:** To retrieve the `RESOURCE_GROUP_GUID` from existing resource groups in your account, 
use the Lacework CLI command `lacework resource-group list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
