terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

variable "name" {
  type = string
  default = "Azure Active Directory Activity Log integration example"
}

resource "lacework_integration_azure_ad_al" "example" {
  name                = var.name
  tenant_id           = "your-tenant-id-goes-here"
  event_hub_namespace = "your-eventhub-ns.servicebus.windows.net"
  event_hub_name      = "your-event-hub-name"
  credentials {
    client_id     = "1234567890-abcd-client-id"
    client_secret = "SUPER_SECURE_SECRET"
  }
  retries = 10
}
