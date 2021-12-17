---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_jira_cloud"
description: |-
  Create and manage Jira Cloud Alert Channel integrations
---

# lacework\_alert\_channel\_jira\_cloud

Configure Lacework to forward events to Jira. Lacework calls the Jira integration REST API and creates a new Jira open issue for each Lacework event that meets or exceeds the specified alert severity level. If there is a large volume of events that exceed the ability of Jira REST API to create new Jira issues, priority is given to those events with the highest severity.

## Jira + Lacework Integration Prerequisites

Before creating the Jira alert channel, verify the following prerequisites:

* Provide a Jira user name and an API access token that is used to create new Jira issues. For management and security purposes, Lacework recommends creating a dedicated Lacework Jira user with appropriate permissions. For more information, refer to the [Jira REST API Reference](https://developer.atlassian.com/server/jira/platform/rest-apis/)
* The Jira user must have sufficient privileges to create new Jira issues in the specified Jira project
* The Jira issue type must exist in the specified Jira project prior to creating the Lacework Jira alert channel. When Lacework creates new Jira issues, it creates new issues based on the specified Jira issue type

## Example Usage

```hcl
resource "lacework_alert_channel_jira_cloud" "example" {
  name        = "My Jira Cloud Alert Channel Example"
  jira_url    = "mycompany.atlassian.net"
  issue_type  = "Bug"
  project_key = "EXAMPLE"
  username    = "my@username.com"
  api_token   = "abcd1234"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `jira_url` - (Required) The URL of your Jira implementation without https protocol (`https://`). For example, `mycompany.atlassian.net` or `mycompany.jira.com`.
* `issue_type` - (Required) The Jira Issue type (such as a `Bug`) to create when a new Jira issue is created.
* `project_key` - (Required) The project key for the Jira project where the new Jira issues should be created.
* `username` - (Required) The Jira user name. Lacework recommends a dedicated Jira user. See above for more information.
* `api_token` - (Required) The Jira API Token. For more information, see [how to create a Jira API Token](https://confluence.atlassian.com/cloud/api-tokens-938839638.html).
* `group_issues_by` - (Optional) Defines how Lacework compliance events get grouped. Must be one of `Events` or `Resources`. Defaults to `Events`.
  The available options are:
  * **Events**:	Single Jira issue will be created when compliance events of the same type but from different resources are detected by Lacework. For example, if three different S3 resources are generating the same compliance event, only one Jira ticket is created.
  * **Resources**: Multiple Jira issues will be created when multiple resources are generating the same compliance event. For example, if three different S3 resources are generating the same compliance event, three Jira issues are created.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `custom_template_file` - (Optional) A Custom Template JSON file to populate fields in the new Jira issues.

## Import

A Lacework Jira Cloud Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_jira_cloud.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).

