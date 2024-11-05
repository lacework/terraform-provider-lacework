---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: convert_to_newer_resource_group"
description: |-
  Converts old to new resource groups
---

# convert_to_newer_resource_group

The new version of Resource Groups changes the Terraform syntax used to create resource groups. 
The old version of Resource Groups defined specific filter fields for each of the Resource Group 
types. The new version defines the `group` argument as an expression tree representing the 
relationships between resources.

Refer to [Filterable Fields section]https://docs.fortinet.com/document/lacework-forticnapp/latest/api-reference/690087/using-the-resource-groups-api
) of the Lacework API documentation for supported resource group filters.

# Lacework Account Resource Group

There is no new version for Lacework account (organization-level) resource groups. No conversion is applicable.

# AWS Resource Groups

## Old Resource Groups

For AWS, old Resource Groups supported the `accounts` list field, which contained the accounts to be included in the resource group. If any of the accounts in the resource group matched an entry in the list, the resource group would be applied.

### Example

```hcl
resource "lacework_resource_group_aws" "example" {
  name        = "My AWS Resource Group"
  description = "A subset of AWS Accounts"
  accounts    = ["123456789", "234567891"]
}
```

## New Resource Groups

New Resource Groups supports an expression tree structure in which the `field` field in a filter object is `Account` and the `value` field contains a list of the account IDs.

### Example

```hcl
resource "lacework_resource_group" "example" {
  name        = "My AWS Resource Group"
  type        = "AWS"
  description = "This groups a subset of AWS resources"
  group {
    operator = "OR"
    filter {
      filter_name = "filter1"
      field     = "Account"
      operation = "EQUALS"
      value     = ["123456789", "234567891"]
    }
  }
}
```

# Azure Resource Groups

## Old Resource Groups

For Azure, the old Resource Groups supported the `tenant` field, which contained a list of subscriptions to be included in the Resource Group.

### Example

```hcl
resource "lacework_resource_group_azure" "example" {
  name          = "My Azure Resource Group"
  description   = "This groups a subset of Azure Subscriptions"
  tenant        = "a11aa1ab-111a-11ab-a000-11aa1111a11a"
  subscriptions = ["1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0", "2b000c3-ab10-1a01-1abc-1a000ab0a0a0"]
}
```

## New Resource Groups

The new Resource Groups supports an expression tree structure in which the `field` field in a filter object is defined as either `Tenant ID` or `Subscription ID` and the `value` field contains a list of the tenant or subscrption IDs.

### Example

```hcl
resource "lacework_resource_group_azure" "example" {
  name        = "My Azure Resource Group"
  type        = "AZURE"
  description = "This groups a subset of Azure Subscriptions"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Tenant ID"
      operation = "EQUALS"
      value     = ["a11aa1ab-111a-11ab-a000-11aa1111a11a"]
    }
    group {
        operator = "OR"
        filter {
          filter_name = "filter2"
          field     = "Subscription ID"
          operation = "EQUALS"
          value     = ["1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0", "2b000c3-ab10-1a01-1abc-1a000ab0a0a0"]
        }
    }
  }
}
```

# GCP Resource Groups

## Old Resource Groups

For GCP, the old Resource Groups supported the `projects` field which contained a list of projects to be included in the resource group.

### Example

```hcl
resource "lacework_resource_group_gcp" "example" {
  name         = "My GCP Resource Group"
  description  = "This groups a subset of Gcp Projects"
  projects     = ["project-1", "project-2", "project-3"]
  organization = "MyGcpOrgID"
}
```

## New Resource Groups

The new Resource Groups supports an expression tree structure in which the `field` field in a filter object is defined as either `Organization ID` or `Project ID` and the `value` field contains a list of the organization or project IDs.

### Example

```hcl
resource "lacework_resource_group_gcp" "example" {
  name        = "My GCP Resource Group"
  type        = "GCP"
  description = "This groups a subset of Gcp Projects"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Organization ID"
      operation = "EQUALS"
      value     = ["MyGcpOrgID"]
    }
    filter {
      filter_name = "filter2"
      field     = "Project ID"
      operation = "EQUALS"
      value     = ["project-1", "project-2", "project-3"]
    }
  }
}
```

# Container Resource Groups

## Old Resource Groups

For containers, the old Resource Groups supported the `container_tags` and `container_label` fields.

### Example

```hcl
resource "lacework_resource_group_container" "example" {
  name           = "My Container Resource Group"
  description    = "This groups a subset of Container Tags"
  container_tags = ["my-container"]
  container_label {
    key   = "name"
    value = "my-container"
  }
}
```

## New Resource Groups

The new Resource Groups supports an expression tree structure in which the `field` field in a filter object is defined as either `Image Tag` or `Container Label` and the `value` field contains a list of the image tags or container labels.

### Example Representation

```hcl
resource "lacework_resource_group_container" "example" {
  name        = "My Container Resource Group"
  type        = "CONTAINER"
  description = "This groups a subset of Container Tags"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Image Tag"
      operation = "EQUALS"
      value     = ["my-container"]
   }
   filter {
      filter_name = "filter2"
      field     = "Container Label"
      operation = "EQUALS"
      value     = ["my-container"]
      key = “name”
   }
  }
}
```

# Machine Resource Groups

## Old Resource Groups

For machines, the old Resource Groups supported the `machine_tags` field.

### Example

```hcl
resource "lacework_resource_group_machine" "example" {
  name        = "My Machine Resource Group"
  description = "This groups a subset of Machine Tags"
  machine_tags {
    key   = "name"
    value = "myMachine"
  }
}
```

## New Resource Groups
 new Resource Groups supports an expression tree structure in which the `field` field in a filter object is defined as either `Image Tag` or `Container Label` and the `value` field contains a list of the image tags or container labels.
 
The new Resource Groups supports an expression tree structure in which the `field` field in a filter object is defined as `Machine Tag` and the `value` field contains a list of the machine tags.

### Example

```hcl
resource "lacework_resource_group_machine" "example" {
  name        = "My Machine Resource Group"
  type        = "MACHINE"
  description = "This groups a subset of Machine Tags"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Machine Tag"
      operation = "EQUALS"
      value     = ["myMachine"]
      key       = “name”
    }
  }
}
```
