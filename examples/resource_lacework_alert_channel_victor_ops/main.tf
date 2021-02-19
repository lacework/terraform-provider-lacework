provider "lacework" {}

resource "lacework_alert_channel_victor_ops" "example" {
  name        = "Victor Ops example"
  webhook_url = "https://alert.victorops.com/integrations/generic/20131114/alert/31e945ee-5cad-44e7-afb0-97c20ea80dd8/database"
}
