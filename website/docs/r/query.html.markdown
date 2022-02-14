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
resource "lacework_query" "MyLQL" {
  id    = "MyLQL"
  query = <<EOT
    {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 's3.amazonaws.com'
      }
      return {
          INSERT_ID
      }
    }
EOT
}
```


## Argument Reference

The following arguments are supported:

* `id` - (Required) The query id.
* `query` - (Required) The query string.
* `evauator_id` - (Optional) The query evaluator id.

## Import

A Lacework query can be imported using a `QUERY_ID`, e.g.

```
$ terraform import lacework_query.example MyLQLQueryID
```
