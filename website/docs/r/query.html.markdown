---
subcategory: "Query"
layout: "lacework"
page_title: "Lacework: lacework_query"
description: |-
  Create and manage Queries
---

# lacework\_query

To provide customizable specification of datasources, Lacework provides the Lacework Query Language (LQL). 
LQL is a human-readable text syntax for specifying selection, filtering, and manipulation of data. 
It overlaps with SQL in its syntax and what it allows.

For more information, see the [LQL Overview Documentation](https://docs.lacework.com/lql-overview).

## Example Usage

```hcl
  resource "lacework_query" "example" {
  query_id       = Lql_Terraform_Query
  query          = <<EOT
    Lql_Terraform_Query {
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }
}
   EOT
}
```


## Argument Reference

The following arguments are supported:

* `query_id` - (Required) The query id.
* `query` - (Required) The query string.
* `evauator_id` - (Optional) The query evaluator id.

## Import

A Lacework query can be imported using a `QUERY_ID`, e.g.

```
$ terraform import lacework_query.example MyLQLQueryID
```

-> **Note:** To retreive the `QUERY_ID` from existing queries in your account, use the
Lacework CLI command `lacework query list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
