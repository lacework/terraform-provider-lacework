provider "lacework" {}

resource "lacework_integration_azure_al" "example" {
  name      = "Azure Activity Log integration example"
  tenant_id = "your-tenant-id-goes-here"
  queue_url = "https://example.queue.core.windows.net/example"
  credentials {
    client_id     = "1234567890-abcd-client-id"
    client_secret = "SUPER_SECURE_SECRET"
  }
}
