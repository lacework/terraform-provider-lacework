---
subcategory: "Queries"
layout: "lacework"
page_title: "Lacework: lacework_query"
description: |-
  Create and manage Lacework Queries
---

# lacework\_query

To provide customizable specification of datasources, Lacework provides the Lacework Query Language (LQL). 
LQL is a human-readable text syntax for specifying selection, filtering, and manipulation of data. 
It overlaps with SQL in its syntax and what it allows.

For more information, see the [LQL Overview Documentation](https://docs.lacework.com/lql-overview).

## Example Usage

Query all EC2 instances with public IP addresses.

```hcl
resource "lacework_query" "example" {
  query_id = "TF_AWS_Config_EC2InstanceWithPublicIPAddress"
  query    = <<EOT
  {
      source {
          LW_CFG_AWS_EC2_INSTANCES
      }
      filter {
          value_exists(RESOURCE_CONFIG:PublicIpAddress)
      }
      return distinct {
          ACCOUNT_ALIAS,
          ACCOUNT_ID,
          ARN as RESOURCE_KEY,
          RESOURCE_REGION,
          RESOURCE_TYPE,
          SERVICE,
          case when RESOURCE_TYPE = 'ec2:instance' then 'HasPublicIp'
          end as COMPLIANCE_FAILURE_REASON
      }
  }
EOT
}
```

Query CloudTrail events and filter only S3 buckets with ACL 'public-read', 'public-read-write' or 'authenticated-read'.

```hcl
resource "lacework_query" "example" {
  query_id       = "TF_AWS_CTA_S3PublicACLCreated"
  query          = <<EOT
  {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 's3.amazonaws.com'
          and EVENT_NAME = 'CreateBucket'
          and EVENT:requestParameters."x-amz-acl"
          in ('public-read', 'public-read-write', 'authenticated-read')
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

## Import

A Lacework query can be imported using a `QUERY_ID`, e.g.

```
$ terraform import lacework_query.example YourLQLQueryID
```

-> **Note:** To retrieve the `QUERY_ID` from existing queries in your account, use the
Lacework CLI command `lacework query list`. To install this tool follow
[this documentation](https://docs.lacework.com/cli/).
