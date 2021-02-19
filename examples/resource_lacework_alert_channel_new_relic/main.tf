provider "lacework" {}

resource "lacework_alert_channel_new_relic" "example" {
  name       = "My New Relic Insights Channel Alert Example"
  account_id = 2338053
  insert_key = "x-xx-xxxxxxxxxxxxxxxxxx"
}
