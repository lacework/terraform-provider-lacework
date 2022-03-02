terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_query" "example" {
  query_id = var.query_id
  query    = var.query
}

variable "query_id" {
  type    = string
  default = "Lql_Terraform_Query"
}

variable "query" {
  type    = string
  default = <<EOT
  {
      source {
          CloudTrailRawEvents
      }
      filter {
          EVENT_SOURCE = 'signin.amazonaws.com'
          and EVENT:userIdentity."type"::String = 'AWSService'
          and EVENT:sourceIPAddress not in ('1.1.1.1', '2.2.2.2')
          and ERROR_CODE is null
      }
      return distinct {
          INSERT_ID, INSERT_TIME, EVENT_TIME, EVENT
      }
  }
EOT
}

output "query_id" {
  value = lacework_query.example.id
}

output "query" {
  value = lacework_query.example.query
}
