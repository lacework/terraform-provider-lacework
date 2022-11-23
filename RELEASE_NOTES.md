# Release Notes
Another day, another release. These are the release notes for the version `v1.1.0`.

## Features
* feat: add event category type 'K8sActivity' to alert_rule resource (#413) (Darren)([47fd9f6](https://github.com/lacework/terraform-provider-lacework/commit/47fd9f6a33a82576bd27ddbf8fca43c445240382))
## Refactor
* refactor: remove last usages of v1 api (#397) (Darren)([724842c](https://github.com/lacework/terraform-provider-lacework/commit/724842cf96aa6d151401e20f678fc742583370f3))
* refactor: remove deprecated fields (#395) (Darren)([65931a3](https://github.com/lacework/terraform-provider-lacework/commit/65931a38cce2e90fd9b2345c0ada300d61e27a41))
* refactor: migrate resource_lacework_integration_docker_v2 to use v2 api (#377) (Darren)([866774f](https://github.com/lacework/terraform-provider-lacework/commit/866774f5abe2183c816b6f15d6f34eeebe4a4228))
* refactor: migrate resource_lacework_integration_docker_hub to use v2 api (#378) (Darren)([0e916c8](https://github.com/lacework/terraform-provider-lacework/commit/0e916c8040fe5a55eb0aef6de3b7422130834646))
* refactor: migrate agent_access_token to use v2 api (#370) (Darren)([917011e](https://github.com/lacework/terraform-provider-lacework/commit/917011e85829c90fe6b34ce7c6b8e8c74fab6a92))
* refactor: migrate resource_lacework_integration_azure_al to use v2 api (#373) (Darren)([d5ca4d4](https://github.com/lacework/terraform-provider-lacework/commit/d5ca4d46c680665bdbcc15850be98cbddc210676))
* refactor: migrate resource_lacework_integration_azure_cfg to use v2 api (#374) (Darren)([8dbf1bf](https://github.com/lacework/terraform-provider-lacework/commit/8dbf1bfbad1cf43523e787c8d9af2b6e2af69f3a))
* refactor: migrate resource_lacework_integration_aws_govcloud_cfg to use v2 api (#375) (Darren)([5d5f082](https://github.com/lacework/terraform-provider-lacework/commit/5d5f082431c7e5b4e9f489fb8250f77c5a67e21f))
* refactor: migrate resource_lacework_integration_aws_govcloud_ct to use v2 api (#376) (Darren)([805d386](https://github.com/lacework/terraform-provider-lacework/commit/805d3867a3795e6ea828c84a9f1fcc21a3a491fd))
## Bug Fixes
* fix(resource): GCR limit_by_label argument (Salim Afiune Maya)([997b861](https://github.com/lacework/terraform-provider-lacework/commit/997b861e2bdbb4514394639811a3a62cebaa542b))
* fix(resource): ECR limit_by_label argument (Salim Afiune Maya)([658d10f](https://github.com/lacework/terraform-provider-lacework/commit/658d10f7503a875265600467fcdce094636f1260))
* fix(resource): docker_v2 limit_by_label argument (Salim Afiune Maya)([37a944e](https://github.com/lacework/terraform-provider-lacework/commit/37a944e61c4d8f609138f7df12980afbe5f19b68))
* fix(resource): docker_hub limit_by_label argument (#416) (Salim Afiune)([91f558c](https://github.com/lacework/terraform-provider-lacework/commit/91f558cdd79bbc6eeb951e7cc51c10bf4e6c9ea9))
* fix: azure_cfg resource sening incorrect cloud account type (#403) (Darren)([1036ecd](https://github.com/lacework/terraform-provider-lacework/commit/1036ecd8ddc7f09b9ebefb7f8ba30d8fa8a02655))
* fix: remove unused func (Darren Murray)([dddeb21](https://github.com/lacework/terraform-provider-lacework/commit/dddeb2157f5d4e6aa997c464e2d91a9fb5a24f60))
* fix: remove unused funcs (Darren Murray)([b6ff8f8](https://github.com/lacework/terraform-provider-lacework/commit/b6ff8f89b88a9f698fa06d4cd2ade5fe7cdece6b))
## Documentation Updates
* docs: add limit_by_label examples (Salim Afiune Maya)([5e8e12c](https://github.com/lacework/terraform-provider-lacework/commit/5e8e12c22dd2c57f761399649126c593735a5200))
* docs: update deprecated cli cmd in docs (#396) (Darren)([1394822](https://github.com/lacework/terraform-provider-lacework/commit/139482264d9c1e5d516b90535c64e2f12a34ca65))
## Other Changes
* build(deps): bump golang.org/x/text from 0.3.7 to 0.4.0 (#388) (dependabot[bot])([83bccc0](https://github.com/lacework/terraform-provider-lacework/commit/83bccc07a080cbba2349432248a62550b8547601))
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#392) (dependabot[bot])([3fc1375](https://github.com/lacework/terraform-provider-lacework/commit/3fc1375cd2e119c1c19de719e41f53827f56a77f))
* build(deps): bump github.com/gruntwork-io/terratest (#386) (dependabot[bot])([eafb1b4](https://github.com/lacework/terraform-provider-lacework/commit/eafb1b4d1e903cb2eec670cacdc3ed9d525c935c))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#387) (dependabot[bot])([5dbbd34](https://github.com/lacework/terraform-provider-lacework/commit/5dbbd34f66e52f5838c9247e2635d8d92749ffea))
* build(deps): bump github.com/lacework/go-sdk from 0.43.0 to 0.44.1 (#393) (dependabot[bot])([43fb505](https://github.com/lacework/terraform-provider-lacework/commit/43fb505b7150a096acee33ce1d11fd9bef7f30e8))
* ci: version bump to v1.0.1-dev (Salim Afiune Maya)([e08bf89](https://github.com/lacework/terraform-provider-lacework/commit/e08bf8995981e08bd16fc7700a0df86f9e338343))
* ci: version bump to v0.27.1-dev (Lacework)([2d8e4f8](https://github.com/lacework/terraform-provider-lacework/commit/2d8e4f882d35ee7d66d09c574e8269de3e75be6f))
* test: change alert_profile test to read element from alerts by name (#412) (Darren)([e9ce33a](https://github.com/lacework/terraform-provider-lacework/commit/e9ce33aba1f665a75abeb75aa2c62fd9b162f003))
* test: update integration test read to use v2 api (Darren Murray)([f5323f6](https://github.com/lacework/terraform-provider-lacework/commit/f5323f604b7be29d4b9ac38c33692a85561c26f0))
* test: move all tests to use APIv2 (#394) (Darren)([b88acee](https://github.com/lacework/terraform-provider-lacework/commit/b88acee715d181e0e20289a4106dd4ad56a7ba86))
