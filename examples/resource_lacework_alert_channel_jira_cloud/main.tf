provider "lacework" {}

resource "lacework_alert_channel_jira_cloud" "example" {
  name        = "My Jira Cloud Alert Channel Example"
  jira_url    = "mycompany.atlassian.net"
  issue_type  = "Bug"
  project_key = "EXAMPLE"
  username    = "my@username.com"
  api_token   = "abcd1234"

  group_issues_by      = "Events"
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
