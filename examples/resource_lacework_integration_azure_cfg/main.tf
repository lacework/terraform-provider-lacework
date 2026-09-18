provider "lacework" {}

resource "lacework_integration_azure_cfg" "example" {
  name      = "Azure config integration example"
  tenant_id = "your-tenant-id-goes-here"
  credentials {
    client_id     = "1234567890-abcd-client-id"
    client_secret = "SUPER_SECURE_SECRET"
  }
}
