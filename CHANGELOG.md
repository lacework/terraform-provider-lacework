## v0.2.13 (February 19, 2021)

## Features
* feat: new IBM QRadar alert channel resource (#76) (Darren)([25e6594](https://github.com/lacework/terraform-provider-lacework/commit/25e65940b55ada83430ccada210ef1e10eb9a87b))
* feat: New Relic Insights alert channel resource (#73) (Darren)([5d7877a](https://github.com/lacework/terraform-provider-lacework/commit/5d7877a9da5155cde1da58de00cbc37cb789f4e9))
* feat: New Victor Ops alert channel resource (#71) (Darren)([641f7cd](https://github.com/lacework/terraform-provider-lacework/commit/641f7cda144eb919abefb6c7b49b93e95384e47c))
* feat: New Cisco Webex alert channel resource (#69) (Darren)([f6e9f8d](https://github.com/lacework/terraform-provider-lacework/commit/f6e9f8d16a026185d142b13546dff8c7f1f286fb))
* feat: new Microsoft Teams alert channel resource (#68) (Darren)([5e4f21f](https://github.com/lacework/terraform-provider-lacework/commit/5e4f21f90fa9a5624bcfc5bb87cffec087f45aa0))
## Documentation Updates
* doc: create service accounts for GCR integration (#74) (Salim Afiune)([be8fb7f](https://github.com/lacework/terraform-provider-lacework/commit/be8fb7fd8eabca92dcefa2eac3ccaab30fdb5deb))
* docs: fix closing curly brace in code snippet for alert channel s3 data export (#70) (Scott Ford)([68c0de6](https://github.com/lacework/terraform-provider-lacework/commit/68c0de6589612cd89454116a46287c3606ba2324))

## v0.2.12 (February 05, 2021)

## Features
* feat: new Datadog alert channel resource (#66) (Darren)([ec90fc4](https://github.com/lacework/terraform-provider-lacework/commit/ec90fc4d98104cb6940d202a33acbf878d4fc83c))

## v0.2.11 (January 28, 2021)

## Features
* feat: new ServiceNow alert channel resource (#63) (Darren)([568a097](https://github.com/lacework/terraform-provider-lacework/commit/568a097b20f1a503361e383645f8e0f8b17db5e1))
## Bug Fixes
* fix: Add issue_grouping field to gcp pub sub (#62) (Darren)([2b64473](https://github.com/lacework/terraform-provider-lacework/commit/2b6447341fe3d7e963ebb0d2a251f62f13608b03))

## v0.2.10 (January 22, 2021)

## Features
* feat: new lacework_alert_channel_gcp_pub_sub resource (#58) (Darren)([a15e6db](https://github.com/lacework/terraform-provider-lacework/commit/a15e6db2f79f00bd5b2bf47c0b35b365d94a87fb))
## Documentation Updates
* docs: fix typo in aws_ct website resource (#60) (Salim Afiune)([815b09b](https://github.com/lacework/terraform-provider-lacework/commit/815b09b7c3ae4d32daa5b71303325c4af0ff7324))

## v0.2.9 (January 15, 2021)

## Features
* feat: new lacework_alert_channel_splunk resource (#56) (Darren)([b1cd97b](https://github.com/lacework/terraform-provider-lacework/commit/b1cd97be62106bdc465bc30b585fe5f10aaae1a8))

## v0.2.8 (January 05, 2021)

## Features
* feat: new lacework_alert_channel_aws_s3 resource (#52) (Darren)([45e39b2](https://github.com/lacework/terraform-provider-lacework/commit/45e39b20acf8ab7f9090b4de230388c20ce3be51))
## Documentation Updates
* docs: update gcp docs for proper render (Salim Afiune Maya)([159595c](https://github.com/lacework/terraform-provider-lacework/commit/159595c01a7bec557c02d561ce670e198b35892b))
* docs: adding more documentation about GCP resources (#53) (Salim Afiune)([ca42d69](https://github.com/lacework/terraform-provider-lacework/commit/ca42d691807ffc1941f48f5c1e77136c34289e13))
## Other Changes
* ci: send slack notifications to team alias ⭐ (Salim Afiune Maya)([887e173](https://github.com/lacework/terraform-provider-lacework/commit/887e17300bd4bdf3863d754e04157d98d060ea2c))

## v0.2.7 (December 17, 2020)

## Features
* feat(resource): add Webhook integration (#49) (Darren)([6f16fc3](https://github.com/lacework/terraform-provider-lacework/commit/6f16fc33e5e8a732540a855b9517d71252bc253e))
## Other Changes
* chore(deps): update go-sdk to remove automatic trigger (#50) (Salim Afiune)([f7af173](https://github.com/lacework/terraform-provider-lacework/commit/f7af1736c0b8ddc9b757bcb73e669783fb2b9409))
* build: upgrade Go version to 1.15 (#48) (Darren)([65cc622](https://github.com/lacework/terraform-provider-lacework/commit/65cc622beabb042a656a22b4a4ef39f930624740))

## v0.2.6 (December 09, 2020)

## Features
* feat(data-source): new lacework_agent_access_token (Salim Afiune Maya)([6283e63](https://github.com/lacework/terraform-provider-lacework/commit/6283e63a3b29df85525cac3044767edf07a97f58))
* feat(resource): new lacework_agent_access_token (#45) (Salim Afiune)([de8b2cc](https://github.com/lacework/terraform-provider-lacework/commit/de8b2cc5c9727a30cd7915e3f3c0032d8b34771d))
* feat(resource): AWS consolidated CloudTrail support (Salim Afiune Maya)([b232e16](https://github.com/lacework/terraform-provider-lacework/commit/b232e163e612b0790b27493df74eb0ac79a840ce))
## Bug Fixes
* fix(docker): avoid loading non-existing password (Salim Afiune Maya)([5ea1e5d](https://github.com/lacework/terraform-provider-lacework/commit/5ea1e5d6d18db86c72baff9ed40151475964d269))
* fix: only encode mapping account when there is one (Salim Afiune Maya)([aedbee8](https://github.com/lacework/terraform-provider-lacework/commit/aedbee8cdf8e2021ecdd14f4daa64795271eb7af))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([ad616a6](https://github.com/lacework/terraform-provider-lacework/commit/ad616a6ef57eb4b19581d5284a660fb1dfe383f5))
* ci: improve release notes generation (Salim Afiune Maya)([11e4604](https://github.com/lacework/terraform-provider-lacework/commit/11e4604e4c6d6919c8f3f005f198e531773f5a0b))

## v0.2.5 (October 13, 2020)

## Features
* feat: trigger initial compliance report automatically (#39) (Salim Afiune)([a9f4dae](https://github.com/lacework/terraform-provider-lacework/commit/a9f4dae0e98a600fa9bdff23c78c578a97cee6d7))
## Other Changes
* ci: zip binaries without bin/ (#38) (Salim Afiune)([2d57b92](https://github.com/lacework/terraform-provider-lacework/commit/2d57b92bbf7f514ed9c11d10a47855d095625fcb))
* ci: add date to CHANGELOG.md (#37) (Salim Afiune)([509e4db](https://github.com/lacework/terraform-provider-lacework/commit/509e4dbd0daa88a63666be1c974554ad14e0b9ce))

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
