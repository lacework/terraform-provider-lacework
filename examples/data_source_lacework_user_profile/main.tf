terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

data "lacework_user_profile" "test" {}

output "lacework_user_profile_url" {
  value = data.lacework_user_profile.test.url
}
