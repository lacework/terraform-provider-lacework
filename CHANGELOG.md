## v0.8.1 (August 18, 2021)

## Bug Fixes
* fix: migrate AWS S3 Alert Channels to API v2 (Salim Afiune Maya)([6fa97ac](https://github.com/lacework/terraform-provider-lacework/commit/6fa97ac001549861526e90dbcf788ae6dab91c9d))
* fix: migrate Slack Alert Channels to use API v2 (Salim Afiune Maya)([d84cc51](https://github.com/lacework/terraform-provider-lacework/commit/d84cc51a1f72e5bd7575f32b1609a1d4766284b0))
## Documentation Updates
* docs: improve documentation for org account_mappings (#146) (Salim Afiune)([edbc4b7](https://github.com/lacework/terraform-provider-lacework/commit/edbc4b7ab0704bb05dd048ce936fe71fec2a6356))
## Other Changes
* chore: fix Go deps (#149) (Salim Afiune)([db496ee](https://github.com/lacework/terraform-provider-lacework/commit/db496eef7423a101b1c58521e97f81a5ba8bd14e))
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([a5f3fbe](https://github.com/lacework/terraform-provider-lacework/commit/a5f3fbe28f511cf1e69e92154b2cb785b40afe30))
* build(deps): bump github.com/gruntwork-io/terratest (#141) (dependabot[bot])([d2d7f2a](https://github.com/lacework/terraform-provider-lacework/commit/d2d7f2a71516cca785fe32b9fcb23070a4a8108c))
* ci: version bump to v0.8.1-dev (Lacework)([5041c63](https://github.com/lacework/terraform-provider-lacework/commit/5041c63f63a74125e7b3c456758d2a1c34ba2bce))

## v0.8.0 (August 06, 2021)

## Features
* feat: Test Alert Channels during TF Create (#133) (Darren)([3df183d](https://github.com/lacework/terraform-provider-lacework/commit/3df183d4f377e2673e9992763ee054da5239ff29))
* feat: new Google Artifact Registry (GAR) resource (#135) (Salim Afiune)([39e0a2a](https://github.com/lacework/terraform-provider-lacework/commit/39e0a2a168e334933ddc6d562f6aaa0c4dfb12cc))
## Refactor
* refactor: test Alert Channels on create and modify (#138) (Salim Afiune)([d58888a](https://github.com/lacework/terraform-provider-lacework/commit/d58888ac7c72c9b6372cb296350c645e93ada5e6))
## Bug Fixes
* fix(lint): missing declaration or govcloud_ct resource (Salim Afiune Maya)([5653228](https://github.com/lacework/terraform-provider-lacework/commit/565322829be14fcf5374beb7b27e120b843b88b6))
* fix(importer): use new Get() v2 func (Salim Afiune Maya)([cd58f1e](https://github.com/lacework/terraform-provider-lacework/commit/cd58f1ecc1becd12b111ae9198c955e8e84e0e86))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk to v0.12.0 (Salim Afiune Maya)([76fa9e7](https://github.com/lacework/terraform-provider-lacework/commit/76fa9e7220aec7b5a00c5ca9e917f987c590559e))
* chore: fixed deps by running make go-vendor (#136) (Salim Afiune)([e359b88](https://github.com/lacework/terraform-provider-lacework/commit/e359b8818db20cffb47b5de114abbc887e94214b))
* build(deps): bump github.com/gruntwork-io/terratest (#130) (dependabot[bot])([4a47a07](https://github.com/lacework/terraform-provider-lacework/commit/4a47a07ab56772064b543ce5d6b13479ced979f3))
* ci: sign lacework-releng commits (#132) (Salim Afiune)([495a264](https://github.com/lacework/terraform-provider-lacework/commit/495a264b4b772ae8d73f0b5184a08cef9b41ead5))
* ci: make clean-test directive (#131) (Salim Afiune)([8d1c9af](https://github.com/lacework/terraform-provider-lacework/commit/8d1c9afa44c9e95e3843756620a5ec5245bff791))
* ci: clean left over files during integration test (#129) (Salim Afiune)([d777278](https://github.com/lacework/terraform-provider-lacework/commit/d77727841592c032c9bf855ed9e9ad6be1c0cab9))

## v0.7.0 (July 27, 2021)

## Features
* feat: new lacework_alert_channel_email resource (Salim Afiune Maya)([0cdb690](https://github.com/lacework/terraform-provider-lacework/commit/0cdb690f7188aac8af9fea7741b7df6509f9badb))
## Bug Fixes
* fix: drifts with lacework_alert_channel_datadog (Salim Afiune Maya)([be38f69](https://github.com/lacework/terraform-provider-lacework/commit/be38f69dfc16f12e0afad782397288310efb37d2))
## Documentation Updates
* docs: update lacework_alert_channel_datadog resource (Salim Afiune Maya)([99cedc7](https://github.com/lacework/terraform-provider-lacework/commit/99cedc701f9f7715c38a243e694b756253634c9c))
* docs: new lacework_alert_channel_email resource (Salim Afiune Maya)([844065a](https://github.com/lacework/terraform-provider-lacework/commit/844065a87a7d7aa96fd0c23427399dc40197366d))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk to v0.11.0 (#126) (Salim Afiune)([5dbefc4](https://github.com/lacework/terraform-provider-lacework/commit/5dbefc44acba06c70902f5e1775af31b1413e022))
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([e282b3e](https://github.com/lacework/terraform-provider-lacework/commit/e282b3ea850a0b5b3aebe333e9032a159122998e))
* ci: fix dependencies (#127) (Salim Afiune)([9cbef84](https://github.com/lacework/terraform-provider-lacework/commit/9cbef84ed82dd54fa68e9b8a9675ed0276ede025))
* ci: fix integration tests (Salim Afiune Maya)([9a6f1e2](https://github.com/lacework/terraform-provider-lacework/commit/9a6f1e2764d10e489fd65a160ad940048932bea2))
* test: e2e testing with terratest (#121) (Darren)([8bc0bd9](https://github.com/lacework/terraform-provider-lacework/commit/8bc0bd9e6b970d2a071698f08d0c0a1f61d2fba7))

## v0.6.0 (July 21, 2021)

## Features
* feat: add provider organization argument to access org data sets (Salim Afiune Maya)([5aa0968](https://github.com/lacework/terraform-provider-lacework/commit/5aa096816bdd03c1ac9386df6a8ff8b5a6374c78))
* feat: add support for multiple tags and labels (#114) (Salim Afiune)([5cf4100](https://github.com/lacework/terraform-provider-lacework/commit/5cf41004b97820bdba90bd263a932daf17e4f724))
* feat: add subaccount arg for org admins (#112) (Salim Afiune)([1d6d7c1](https://github.com/lacework/terraform-provider-lacework/commit/1d6d7c1fc0e36f3179d2b7b428475819290c69dc))
## Refactor
* refactor: use V2 CloudAccounts for integration_aws_ct resource (Salim Afiune Maya)([0c1cc05](https://github.com/lacework/terraform-provider-lacework/commit/0c1cc05b1c29c5abc81a872ba2441f51652cb50d))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([e4aa1da](https://github.com/lacework/terraform-provider-lacework/commit/e4aa1daf365a433795b0d2eccb6fd0db7ea07b66))
* build(deps): bump github.com/lacework/go-sdk to v0.10.1 (#122) (Salim Afiune)([0d158ce](https://github.com/lacework/terraform-provider-lacework/commit/0d158cedfc9fa966c3939360c2c06c1c2e0a2b69))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#120) (dependabot[bot])([c1201ba](https://github.com/lacework/terraform-provider-lacework/commit/c1201ba44201f031470507ee6ccf1fab7c487b1a))
* build(deps): bump github.com/lacework/go-sdk from 0.8.0 to 0.9.1 (#117) (dependabot[bot])([23bb9ce](https://github.com/lacework/terraform-provider-lacework/commit/23bb9ce2cdb9615b1d93fc143cb51574aee4be57))
* ci: pin Go tooling to specific versions (#115) (Salim Afiune)([f3db1cc](https://github.com/lacework/terraform-provider-lacework/commit/f3db1cc3b628e56df4a9f857c3a12b86a75f61fe))

## v0.5.0 (May 27, 2021)

## Features
* feat: New AWS CouldTrail for AWS GovCloud resource (#108) (Darren)([17ff48f](https://github.com/lacework/terraform-provider-lacework/commit/17ff48fcdc982497391b6d6f48c186def3c19185))
* feat: New AWS Config for AWS GovCloud resource (#107) (Darren)([d5e0d2e](https://github.com/lacework/terraform-provider-lacework/commit/d5e0d2e0d4e04e77d8c1237d65d90f1446d29be0))
## Documentation Updates
* docs: update GovCloud registry documentation (#109) (Salim Afiune)([1efa3a8](https://github.com/lacework/terraform-provider-lacework/commit/1efa3a8563364806a8a1f5f320f0556348c6a3b8))
## Other Changes
* ci: fix prepare-release job (#110) (Salim Afiune)([c69c3ec](https://github.com/lacework/terraform-provider-lacework/commit/c69c3ec2644a12e0cb06e7bbb02880c16d85f1e4))

## v0.4.1 (May 18, 2021)

## Other Changes
* build: Update go version to 1.16 (#102) (Darren)([a38f559](https://github.com/lacework/terraform-provider-lacework/commit/a38f55961c89bf61a02c1c1547058fe79bafce45))
* ci: fix go deps (#104) (Darren)([3e557ae](https://github.com/lacework/terraform-provider-lacework/commit/3e557ae8cdafe30ed203734eb8bed9fb03000daa))

## v0.4.0 (May 12, 2021)

## Features
* feat: understand Lacework CLI config v2 (Salim Afiune Maya)([49ba175](https://github.com/lacework/terraform-provider-lacework/commit/49ba175a3ba1b8d06cea387460e2b5c715360935))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([ac6f892](https://github.com/lacework/terraform-provider-lacework/commit/ac6f892e1de412c01c5bebc1476aac424a74facc))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#96) (dependabot[bot])([84796e9](https://github.com/lacework/terraform-provider-lacework/commit/84796e9e6041dd73b282a8a721f3da60a76868a7))
* ci: run make prepare (#100) (Salim Afiune)([349d26f](https://github.com/lacework/terraform-provider-lacework/commit/349d26f4949a4aae32f411e7cb7810262bb1cd9c))
* ci: remove CircleCI completely (Salim Afiune Maya)([1fd808a](https://github.com/lacework/terraform-provider-lacework/commit/1fd808aa30f76d0a6d3f2671d2d531fad0582a11))
* ci: update release.sh script to work with CodeFresh (Salim Afiune Maya)([bec5242](https://github.com/lacework/terraform-provider-lacework/commit/bec52425b537a377228528fbb9729d45d08f9ed5))

## v0.3.2 (April 28, 2021)

## Other Changes
* build(deps): bump github.com/lacework/go-sdk from 0.4.0 to 0.6.0 (#95) (dependabot[bot])([b66ceca](https://github.com/lacework/terraform-provider-lacework/commit/b66ceca787313251040d1ae1508c3057116d67fe))
* build(deps): bump github.com/lacework/go-sdk from 0.2.22 to 0.4.0 (#94) (dependabot[bot])([d7c4852](https://github.com/lacework/terraform-provider-lacework/commit/d7c4852927bab0e152b614a9f0a427f4f7889ecf))

## v0.3.1 (March 16, 2021)

## Documentation Updates
* docs(website): fix ECR Module example (#88) (Salim Afiune)([f0a02f2](https://github.com/lacework/terraform-provider-lacework/commit/f0a02f2587a451ae1efad3040d8389f4265ea9d5))

## v0.3.0 (March 16, 2021)

## Features
* feat: lacework_integration_ecr support for IAM Roles (Salim Afiune Maya)([42bfbae](https://github.com/lacework/terraform-provider-lacework/commit/42bfbae107390e959a6684e6322e25ba37e42777))
## Bug Fixes
* fix: lacework_integration_ecr from SCHEMA changes (Salim Afiune Maya)([5a52b98](https://github.com/lacework/terraform-provider-lacework/commit/5a52b98823d1fb1c8942d4c9c8bbaafa952e4eab))
## Documentation Updates
* docs: update ECR argument descriptions (Salim Afiune Maya)([dea6d49](https://github.com/lacework/terraform-provider-lacework/commit/dea6d49019047eb925b2440c0cad1cf46e2c71f0))
* docs(website): update lacework_integration_ecr resource (Salim Afiune Maya)([3e07cd4](https://github.com/lacework/terraform-provider-lacework/commit/3e07cd47b469189e159bdca2c7514b52d6bbddd3))
## Other Changes
* chore(deps): update github.com/lacework/go-sdk (Salim Afiune Maya)([987e527](https://github.com/lacework/terraform-provider-lacework/commit/987e527b8263a519c32c8ca215b8e50cec8cbb30))
* ci: increase timeout of the build step to 30m (Salim Afiune Maya)([b86e6c7](https://github.com/lacework/terraform-provider-lacework/commit/b86e6c7c9e313165a8050020995145488e057c7d))

## v0.2.14 (March 08, 2021)

## Features
* feat: retry mechanism for all cloud resources (#81) (Salim Afiune)([61893af](https://github.com/lacework/terraform-provider-lacework/commit/61893afaedeeb407e5b9f20ef7678b64c37965f4))
* feat: Add support for ServiceNow Alert custom JSON template (#78) (Darren)([a315965](https://github.com/lacework/terraform-provider-lacework/commit/a3159658bda69fc7b54bde758734af51151a6e0a))
## Documentation Updates
* docs(website): add missing retries argument (#83) (Salim Afiune)([0df8c10](https://github.com/lacework/terraform-provider-lacework/commit/0df8c10da7f108bf5c5d669a7eada3906add7391))
* docs: Add custom_template_file field to Service Now docs (Darren Murray)([428aaad](https://github.com/lacework/terraform-provider-lacework/commit/428aaad94aa2210243d43e65a32a4ab2981d8608))
## Other Changes
* chore: Update go-sdk dependencies (Darren Murray)([d0ebfc5](https://github.com/lacework/terraform-provider-lacework/commit/d0ebfc580025e03c03190325edd2ce9cdea1995e))
* chore: Update terraform-plugin-sdk to v2 (#79) (Darren)([ec8299f](https://github.com/lacework/terraform-provider-lacework/commit/ec8299fcee941efdc0020dae78889ad7656c2752))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#80) (dependabot[bot])([d3f6fd7](https://github.com/lacework/terraform-provider-lacework/commit/d3f6fd71bbf8dbdc9c0026c5e2dcaba7a42e2017))
* build(deps): bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (#75) (dependabot[bot])([70204ea](https://github.com/lacework/terraform-provider-lacework/commit/70204eaf6aa64bd9e76799f8a03c3110c637ff2f))
* ci: delete 'master' branch (#84) (Salim Afiune)([8012b0b](https://github.com/lacework/terraform-provider-lacework/commit/8012b0b49bb6acf7743abdcc1614e332a6dda14c))

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
