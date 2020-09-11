## 0.3.0 (Unreleased)
## 0.2.4 (September 11, 2020)

## New Resources

Management of Container Registry integrations.
* `lacework_integration_docker_hub`: Docker Hub integration
* `lacework_integration_docker_v2`: Docker V2 Registry integration
* `lacework_integration_ecr`: Amazon Container Registry (ECR) integration
* `lacework_integration_gcr`: Google Container Registry (GCR) integration

## Other Changes
* docs(website): update all integration names (Salim Afiune Maya)([13d5259](https://github.com/terraform-providers/terraform-provider-lacework/commit/13d525987ddba9841eb2781b31ca8284911a72f1))
* docs(website): change provider sidebar by grouping (Salim Afiune Maya)([3ca830f](https://github.com/terraform-providers/terraform-provider-lacework/commit/3ca830f04e0fe99bceef169d0b0c5240df8d74dc))
* docs: update Lacework CLI markdown link (Salim Afiune Maya)([94d7ea8](https://github.com/terraform-providers/terraform-provider-lacework/commit/94d7ea852342c1ef4927cd64e154178b3e1edfa0))
* chore: use integration ENUM in log messages (Salim Afiune Maya)([c98ed22](https://github.com/terraform-providers/terraform-provider-lacework/commit/c98ed226a1e60a0ae05a8aef0dde1b0b74c393d0))
* chore: ran terraform fmt for all examples/ (Salim Afiune Maya)([7e5dd1b](https://github.com/terraform-providers/terraform-provider-lacework/commit/7e5dd1b5ccd4a1533d4418f1c0490467990b7972))
* fix(azure): suppress client_secret diff to avoid updates (#25) (Salim Afiune)([6fe5dc5](https://github.com/terraform-providers/terraform-provider-lacework/commit/6fe5dc5d0051cc13e2d00492b607db72fb905b9e))

## 0.2.3 (August 27, 2020)

Both, `lacework_alert_channel_jira_cloud` and `lacework_alert_channel_jira_server` have now
a new optional argument named `custom_template_file` to specify a Custom Template file in JSON
format to populate fields in the new Jira issues.

## 0.2.2 (August 14, 2020)

## New Resources

Management of alert channel integrations:
* `lacework_alert_channel_jira_cloud`: Create and manage Jira Cloud Alert Channel integrations
* `lacework_alert_channel_jira_server`: Create and manage Jira Server Alert Channel integrations

## 0.2.1 (July 24, 2020)

## Bug Fixes
* fix(gcp): detect a change of credentials [#14](https://github.com/terraform-providers/terraform-provider-lacework/pull/14)

## 0.2.0 (July 23, 2020)

## Improvements

Support of additional methods to configure the Lacework provider:
* Static credentials
* Environment variables
* Configuration file

## New Resources

Management of alert channel integrations:
* `lacework_alert_channel_aws_cloudwatch`: Forward alerts to an AWS CloudWatch event bus
* `lacework_alert_channel_pagerduty`: Forward alerts to PagerDuty
* `lacework_alert_channel_slack`: Forward alerts to a Slack channel

## Bug Fixes

* Unable to read inputs from provider [#6](https://github.com/terraform-providers/terraform-provider-lacework/issues/6)
* Avoid updating GCP resources on second `terraform apply` [#4](https://github.com/terraform-providers/terraform-provider-lacework/issues/4)

# 0.1.1 (July 09, 2020)

## Improvements
* Add `User-Agent` header for Lacework API backend metrics

# 0.1.0 (June 02, 2020)

## New Resources

Management of external integrations.
* `lacework_integration_aws_cfg`: AWS configuration integration
* `lacework_integration_aws_ct`:  AWS CloudTrail integration
* `lacework_integration_azure_cfg`: Azure configuratio integration
* `lacework_integration_azure_al`: Azure Activity Log integration
* `lacework_integration_gcp_cfg`: GCP configuration integration
* `lacework_integration_gcp_at`: GCP Audit Log integration

Data source to integrate with the Lacework platform.
* `lacework_api_token`: generates a new Lacework API token

### Importers
All resources have the ability to be imported, for existing
Lacework integrations use the import command:
```bash
$ terraform import lacework_integration_aws_cfg.name <INT_GUID>
```
