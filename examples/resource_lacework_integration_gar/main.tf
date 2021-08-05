terraform {
  required_providers {
    lacework = {
      source = "lacework/lacework"
    }
  }
}

provider "lacework" {}

resource "lacework_integration_gar" "example" {
  name            = "Google Artifact Registry Example"
  registry_domain = "us-west1-docker.pkg.dev"
  credentials {
    client_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    private_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    client_email   = "email@some-project-name.iam.gserviceaccount.com"
    private_key    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }

  limit_num_imgs        = 10
  limit_by_tags         = ["dev*", "*test"]
  limit_by_repositories = ["repo/my-image", "repo/other-image"]

  limit_by_label {
    key   = "key"
    value = "value"
  }

  limit_by_label {
    key   = "key"
    value = "value2"
  }

  limit_by_label {
    key   = "foo"
    value = "bar"
  }
}
