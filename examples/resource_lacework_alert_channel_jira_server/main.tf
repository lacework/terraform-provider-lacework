provider "lacework" {}

resource "lacework_alert_channel_jira_server" "example" {
  name        = "My Jira Server Alert Channel Example"
  jira_url    = "mycompany.atlassian.net"
  issue_type  = "Bug"
  project_key = "EXAMPLE"
  username    = "my@username.com"
  password    = "my-password"

  group_issues_by      = "Resources"
  custom_template_file = <<TEMPLATE
{
    "fields": {
        "labels": [
            "myLabel"
        ],
        "priority":
        {
            "id": "1"
        }
    }
}
TEMPLATE
}
