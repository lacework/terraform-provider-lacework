terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

  resource "lacework_query" "example" {
    id    = var.query_id
    query = <<EOT
    MyLQL {
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

variable "query_id" {
  type = string
  default = "lql-terraform-query"
}

output "query_id" {
  value = lacework_query.example.id
}

output "query" {
  value = lacework_query.example.query
}