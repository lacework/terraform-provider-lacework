terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

resource "lacework_query" "example" {
  query = var.query
}

variable "query" {
  type    = string
  default = <<EOT
  Lql_Terraform_Query {
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
