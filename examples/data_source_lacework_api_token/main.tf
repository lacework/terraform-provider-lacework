provider "lacework" {}

data "lacework_api_token" "test" {}

output "lacework_api_token" {
  value = data.lacework_api_token.test.token
}
