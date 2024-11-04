---
subcategory: "Resource Groups"
layout: "lacework"
page_title: "Lacework: convert_to_newer_resource_group"
description: |-
  Converts V1 to V2 resource groups
---

# convert_to_newer_resource_group

The new version of resource groups changes the Terraform syntax used to create the resource groups. 
The older version of resource groups defined specific filter fields for each of the resource group 
types. The newer version defines a group argument for an expression tree representing the 
relationships between resources.

Refer to Filterable Fields section [here](https://docs.fortinet.com/document/lacework-forticnapp/latest/api-reference/690087/using-the-resource-groups-api
) for supported resource group filters

# Lacework Account Resource Groups

There is no new version for lacework account resource groups so no conversion is applicable for 
these organizational level resource groups.

# AWS Resource Groups

## Older Resource Group

The older AWS resource group supported the accounts list field to consist of the accounts to be
included in the resource group. If any of the accounts in the resource group matched an entry in the
list, the resource group would be applied.

### Example Representation

```hcl
resource "lacework_resource_group_aws" "example" {
      name        = "My AWS Resource Group"
      description = "This groups a subset of AWS Accounts"
      accounts    = ["123456789", "234567891"]
}
```

## Newer Resource Group

The new resource group supports an expression tree structure in which the field key in the filter
object would be called Account and the value key would be a list containing the values.

### Example Representation

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

## Older Resource Group

The older Azure resource group supported the tenant field and a list of subscriptions within that 
tenant to be included in the resource group.

### Example Representation

```hcl
resource "lacework_resource_group_azure" "example" {
  name          = "My Azure Resource Group"
  description   = "This groups a subset of Azure Subscriptions"
  tenant        = "a11aa1ab-111a-11ab-a000-11aa1111a11a"
  subscriptions = ["1a1a0b2-abc0-1ab1-1abc-1a000ab0a0a0", "2b000c3-ab10-1a01-1abc-1a000ab0a0a0"]
}
```

## Newer Resource Group

The new resource group representation for the resource group above would support an expression 
tree structure in which the field key in the filter objects would be called either Tenant ID or Subscription ID and the value key would be a list 
containing the values.

### Example Representation

```hcl
resource "lacework_resource_group_azure" "example" {
  name        = "My Azure Resource Group"
  type        = "AZURE"
  description = "This groups a subset of Azure Subscriptions"
  group {
    operator = "AND"
    filter {
      filter_name = "filter1"
      field     = "Tenant"
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

## Older Resource Group

The older GCP resource group supported the organization field and a list of projects within that 
organization to be included in the resource group.

### Example Representation

```hcl
resource "lacework_resource_group_gcp" "example" {
  name         = "My GCP Resource Group"
  description  = "This groups a subset of Gcp Projects"
  projects     = ["project-1", "project-2", "project-3"]
  organization = "MyGcpOrgID"
}
```

## Newer Resource Group

The new resource group representation for the resource group above would support an expression
tree structure in which the field key in the filter objects would be called either Organization 
ID or Project ID and the value key would be a list containing the values.

### Example Representation

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

## Older Resource Group

The older Container resource group supported the fields container tags and labels.

### Example Representation

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

## Newer Resource Group

The new resource group representation for the resource group above would support an expression
tree structure in which the field key in the filter objects would be called either Image Tag or 
Container Label and the value key would be a list containing the values.

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

## Older Resource Group

The older Machine resource group supported the machine tags field.

### Example Representation

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

## Newer Resource Group

The new resource group representation for the resource group above would support an expression
tree structure in which the field key in the filter objects would be called Machine Tag and the 
value key would be a list containing the values.

### Example Representation

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
