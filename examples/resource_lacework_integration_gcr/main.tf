terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_gcr" "example" {
  name            = "GRC Example"
  registry_domain = "gcr.io"
  non_os_package_support = true
  credentials {
    client_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    private_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_email   = "email@some-project-name.iam.gserviceaccount.com"
    private_key    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["my-repo", "other-repo"]

  limit_by_labels = {
    key1 = "label1"
    key2 = "label2"
  }
}
