## 0.2.0 (Unreleased)

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
