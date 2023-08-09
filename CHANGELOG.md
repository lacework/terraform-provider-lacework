## v1.12.0 (August 09, 2023)

## Features
* feat: support new properties in alert_rules (#518) (Darren)([24091a7b](https://github.com/lacework/terraform-provider-lacework/commit/24091a7be364679d7e08c00354f0fe1d19a13b35))
## Bug Fixes
* fix: alert rules remove omitempty (#528) (jonathan stewart)([2f012178](https://github.com/lacework/terraform-provider-lacework/commit/2f0121785afe78dcddf8015241316ed369ec18fc))
* fix: lint error in resource_lacework_managed_policies.go (#525) (Pengyuan Zhao)([cd659b8e](https://github.com/lacework/terraform-provider-lacework/commit/cd659b8efd517c1af410ba4250362dc607ed3dd5))
## Other Changes
* chore: 🧹 clean a bunch of deprecated code (#522) (Salim Afiune)([67ee493d](https://github.com/lacework/terraform-provider-lacework/commit/67ee493dcaf496c2cd5d26306fd4f4aba724d443))
* ci: version bump to v1.10.1-dev (Lacework)([5f5a44f4](https://github.com/lacework/terraform-provider-lacework/commit/5f5a44f42990d5dc1a62e29d14f1e2e5ccabce12))

## v1.11.0 (August 08, 2023)

## Features
* feat: support new properties in alert_rules (#518) (Darren)([24091a7b](https://github.com/lacework/terraform-provider-lacework/commit/24091a7be364679d7e08c00354f0fe1d19a13b35))
## Bug Fixes
* fix: alert rules remove omitempty (#528) (jonathan stewart)([2f012178](https://github.com/lacework/terraform-provider-lacework/commit/2f0121785afe78dcddf8015241316ed369ec18fc))
* fix: lint error in resource_lacework_managed_policies.go (#525) (Pengyuan Zhao)([cd659b8e](https://github.com/lacework/terraform-provider-lacework/commit/cd659b8efd517c1af410ba4250362dc607ed3dd5))
## Other Changes
* chore: 🧹 clean a bunch of deprecated code (#522) (Salim Afiune)([67ee493d](https://github.com/lacework/terraform-provider-lacework/commit/67ee493dcaf496c2cd5d26306fd4f4aba724d443))
* ci: version bump to v1.10.1-dev (Lacework)([5f5a44f4](https://github.com/lacework/terraform-provider-lacework/commit/5f5a44f42990d5dc1a62e29d14f1e2e5ccabce12))

## v1.10.0 (July 25, 2023)

## Features
* feat: add new lacework_managed_policies resource (#516) (Pengyuan Zhao)([5d4f4957](https://github.com/lacework/terraform-provider-lacework/commit/5d4f495734cb1fa0e77fc7e362beeac6a3f1c0fc))
## Documentation Updates
* docs: add documentation for the OCI config integration (#509) (Kolbeinn)([c0ac0c81](https://github.com/lacework/terraform-provider-lacework/commit/c0ac0c813b3d86c8eac5d2e570a9aed9dac93359))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#513) (dependabot[bot])([e4fd1eec](https://github.com/lacework/terraform-provider-lacework/commit/e4fd1eec7e463b414fdf40bf38d2a0dee2c167d4))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#504) (dependabot[bot])([7969b8fa](https://github.com/lacework/terraform-provider-lacework/commit/7969b8fa6afa04b39fc4b64c68f82f83a891f76c))
* build(deps): bump github.com/gruntwork-io/terratest (#503) (dependabot[bot])([e99f3744](https://github.com/lacework/terraform-provider-lacework/commit/e99f3744350aacc6e49c1fb4c9cc82c20785993d))
* build(deps): bump github.com/lacework/go-sdk from 1.24.0 to 1.26.0 (#507) (dependabot[bot])([cbaeea70](https://github.com/lacework/terraform-provider-lacework/commit/cbaeea70d8f65e7c2af7598f9aea1e273ff5c4cf))
* build(deps): bump google.golang.org/grpc from 1.51.0 to 1.53.0 (#506) (dependabot[bot])([4b757e81](https://github.com/lacework/terraform-provider-lacework/commit/4b757e8126244ba18a72fc98ac6359ab7f816f53))
* ci: version bump to v1.9.1-dev (Lacework)([32d757ba](https://github.com/lacework/terraform-provider-lacework/commit/32d757ba41295a9b1fc83631a4915d1c1d055eed))

## v1.9.0 (July 05, 2023)

## Features
* feat(client-agentless): add multi-volume and scan stopped instances fields to agentless integrations (Whitney Smith)([cd864e50](https://github.com/lacework/terraform-provider-lacework/commit/cd864e50061d3f12e400d3605a298cb8f76a4530))
* feat: add OCI integration resource (Kolbeinn Karlsson)([5a361e54](https://github.com/lacework/terraform-provider-lacework/commit/5a361e54180a11a70806f41eb2b3c45881230540))
## Documentation Updates
* docs(website): add multivolume and stopped instance arguments (#505) (Joseph Wilder)([b800edf8](https://github.com/lacework/terraform-provider-lacework/commit/b800edf8072589c1d3f2c37e82599505851d05bb))
## Other Changes
* ci: version bump to v1.8.1-dev (Lacework)([b8011be1](https://github.com/lacework/terraform-provider-lacework/commit/b8011be1f3e961b692a44aca9fd4ac48367a122f))

## v1.8.0 (June 08, 2023)

## Features
* feat(resource): add support for compliance policies (#492) (Pengyuan Zhao)([68b5cd84](https://github.com/lacework/terraform-provider-lacework/commit/68b5cd840761a49c349d16ccee460cd7b66ac28a))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#494) (dependabot[bot])([b4bd9481](https://github.com/lacework/terraform-provider-lacework/commit/b4bd9481e77c18a3f41186b202d703bc2dbf82be))
* build(deps): bump github.com/stretchr/testify from 1.8.3 to 1.8.4 (#495) (dependabot[bot])([b6ce6003](https://github.com/lacework/terraform-provider-lacework/commit/b6ce60032e1a113a18cdc11b710a2685ee9f666c))
* ci: version bump to v1.7.1-dev (Lacework)([57dc1261](https://github.com/lacework/terraform-provider-lacework/commit/57dc126134488d8eefecb9fbc710f3bee838a1e1))

## v1.7.0 (June 02, 2023)

## Features
* feat: Add org account mappings to Agentless for AWS (#473) (Whitney Smith)([a80f2d3f](https://github.com/lacework/terraform-provider-lacework/commit/a80f2d3f99273ef9e53b66c2d42f9bc6a0cb6b91))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#488) (dependabot[bot])([f235c8d8](https://github.com/lacework/terraform-provider-lacework/commit/f235c8d87efbe1a7d019aefe6e5c62645e10def7))
* build(deps): bump github.com/stretchr/testify from 1.8.2 to 1.8.3 (#487) (dependabot[bot])([c106c2be](https://github.com/lacework/terraform-provider-lacework/commit/c106c2bede9394b1c3a73d4f6e53bd2c1fda4315))
* ci: version bump to v1.6.3-dev (Lacework)([9eac643c](https://github.com/lacework/terraform-provider-lacework/commit/9eac643c2b0b12220d69cac80967799212a82eb7))

## v1.6.2 (May 05, 2023)

## Bug Fixes
* fix: update go-sdk version for proxy settings fix (#480) (Darren)([a13f2bde](https://github.com/lacework/terraform-provider-lacework/commit/a13f2bde5952c3a6646f49cb79adc386dbd4dd0c))
## Other Changes
* ci: version bump to v1.6.2-dev (Lacework)([4f534c0d](https://github.com/lacework/terraform-provider-lacework/commit/4f534c0db9441207edba254cfe5c8892813d95b2))

## v1.6.1 (May 02, 2023)

## Other Changes
* build(deps): bump github.com/lacework/go-sdk (#472) (dependabot[bot])([bc40e04c](https://github.com/lacework/terraform-provider-lacework/commit/bc40e04ce8a5bb8c6046434457d88774853986d1))
* build(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (#471) (dependabot[bot])([3b650e81](https://github.com/lacework/terraform-provider-lacework/commit/3b650e812c687b3fa215b5e3a27dcb65bdf2bb23))
* build(deps): bump github.com/gruntwork-io/terratest (#468) (dependabot[bot])([412657e6](https://github.com/lacework/terraform-provider-lacework/commit/412657e68e802a3999ac5b94fe1e43c398a664bf))
* build(deps): bump golang.org/x/text from 0.4.0 to 0.9.0 (#465) (dependabot[bot])([070c1012](https://github.com/lacework/terraform-provider-lacework/commit/070c101275bc098b0b987e85f11447a758cf744e))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#462) (dependabot[bot])([3b756a51](https://github.com/lacework/terraform-provider-lacework/commit/3b756a514b7653d93ec4e8f79a9e7b59677420bc))
* build(deps): bump github.com/hashicorp/go-getter from 1.6.1 to 1.7.0 (#447) (dependabot[bot])([1b673bd8](https://github.com/lacework/terraform-provider-lacework/commit/1b673bd816b7be2f58358441dc849ccd226e8f3c))
* ci: version bump to v1.6.1-dev (Lacework)([97cb3120](https://github.com/lacework/terraform-provider-lacework/commit/97cb312020d1e0005e23cbcc3723d82f4b3abca6))

## v1.6.0 (April 14, 2023)

## Features
* feat: Addition of S3 Bucket ARN for AWS EKS Audit log (#467) (djmctavish)([09309d32](https://github.com/lacework/terraform-provider-lacework/commit/09309d3232d351e64e89f340969e6155baa46b6d))
## Documentation Updates
* docs: add 'enabled' to agent_access_token doc (#466) (Darren)([ea55405c](https://github.com/lacework/terraform-provider-lacework/commit/ea55405cd827a8b6c68730b0abc04dea366ad764))
## Other Changes
* ci: fix release (#469) (Darren)([cb54378f](https://github.com/lacework/terraform-provider-lacework/commit/cb54378f4915ecbd71ce6304e5051564fee49041))
* ci: version bump to v1.5.2-dev (Lacework)([7f954384](https://github.com/lacework/terraform-provider-lacework/commit/7f954384a7c07796de74a1b339370e863ed3f1dd))

## v1.5.1 (March 23, 2023)

## Bug Fixes
* fix: avoid sending an empty GCP private_key_id (#461) (Salim Afiune)([c0efc07b](https://github.com/lacework/terraform-provider-lacework/commit/c0efc07b40ed01df757156e12268b773b31b62e8))
## Other Changes
* chore: add growth-team to CODEOWNERS file (#454) (Darren)([c6e59608](https://github.com/lacework/terraform-provider-lacework/commit/c6e59608ef12b32d7fb49ea33f7c4fd9bc39d85a))
* ci: set gcp private key to env var (#458) (djmctavish)([c6ee536a](https://github.com/lacework/terraform-provider-lacework/commit/c6ee536a686ea26df0308d70c27af3fa5ffb3b81))
* ci: correct gcp private key env var name (#457) (Darren)([1859efbc](https://github.com/lacework/terraform-provider-lacework/commit/1859efbc5411304ec6a0de2fc4036f992ade428a))
* ci: set gcp private key to env var (#456) (Darren)([6a4b8d33](https://github.com/lacework/terraform-provider-lacework/commit/6a4b8d3362fef1b925b77dfd9526fa6691d1b4e9))
* ci: version bump to v1.5.1-dev (Lacework)([b6dcb8f3](https://github.com/lacework/terraform-provider-lacework/commit/b6dcb8f3bc55df9387e089a34abab6f1bedc8753))

## v1.5.0 (March 06, 2023)

## Features
* feat: Addition of support for new Gcp pub sub audit log integration (#444) (djmctavish)([a1d79277](https://github.com/lacework/terraform-provider-lacework/commit/a1d792772c489bda4ea2993c30feacbf7f364f8e))
## Bug Fixes
* fix: alert-rule 'event categories' drift (#446) (Darren)([51fad261](https://github.com/lacework/terraform-provider-lacework/commit/51fad2614099acea14cebb7235bd81f0945a2c8e))
## Other Changes
* ci: disable Team Member integration tests (#441) (Salim Afiune)([dfa28268](https://github.com/lacework/terraform-provider-lacework/commit/dfa28268af3c43ac322fc4925955746235f84e39))
* ci: fix new policy tag restrictions (#437) (Salim Afiune)([afb8c387](https://github.com/lacework/terraform-provider-lacework/commit/afb8c387a750ac5f7ee9ddc1b4897a45b1b2ad91))
* ci: version bump to v1.4.1-dev (Lacework)([ea2d4712](https://github.com/lacework/terraform-provider-lacework/commit/ea2d47124358eb3fce1fa4e45968a86bf257e642))
* test: add new integration test for invalid policy exception constraints (#451) (Darren)([4c7ad874](https://github.com/lacework/terraform-provider-lacework/commit/4c7ad8741ab8e05d99e1cb5d0d80e7ff73bf8cc1))

## v1.4.0 (December 19, 2022)

## Features
* feat: add inline and proxy scanner resources (#434) (Salim Afiune)([d7e00749](https://github.com/lacework/terraform-provider-lacework/commit/d7e00749f4aaeaeb7b80d8741ed9640d331a0e34))
## Other Changes
* ci: version bump to v1.3.1-dev (Lacework)([2cf43092](https://github.com/lacework/terraform-provider-lacework/commit/2cf4309254046cf683bc0971edc9b0ccbbbc7df5))

## v1.3.0 (December 13, 2022)

## Features
* feat(jira): add bidirectional configuration (hazedav)([f8de90c](https://github.com/lacework/terraform-provider-lacework/commit/f8de90cbeb188fd2dfb0e99034df852182d21601))
## Bug Fixes
* fix: add server token and uri computed fields to lacework_gcp_agentless_scanning (#427) (ammarekbote)([cc24dd9](https://github.com/lacework/terraform-provider-lacework/commit/cc24dd96d26c42d5fd1f31184ad0b0e488bf1681))
## Other Changes
* ci: version bump to v1.2.1-dev (Lacework)([8348f4c](https://github.com/lacework/terraform-provider-lacework/commit/8348f4c0721ea2938f007bcda343fe7c79fa539d))

## v1.2.0 (December 13, 2022)

## Features
* feat(api) Adds GCP Sidekick Integration support (Ammar Ekbote)([372a7d0](https://github.com/lacework/terraform-provider-lacework/commit/372a7d0431d39d1245e88b73d77c84ae9bf7e1cd))
## Bug Fixes
* fix: upgrade issue with limit_by_labels (#428) (Darren)([6c06a3c](https://github.com/lacework/terraform-provider-lacework/commit/6c06a3c94d9c4b285b4375a5208f6b0d8ee85f92))
* fix: upgrade issue with limit_by_labels (#426) (Darren)([1e7b662](https://github.com/lacework/terraform-provider-lacework/commit/1e7b6628c155966a9222687d92f993ac7317038d))
## Documentation Updates
* docs: include K8sActivity event category type in alert_rule docs (#414) (Darren)([41af3a7](https://github.com/lacework/terraform-provider-lacework/commit/41af3a7e5cfc5830b30670e67fe2878032abe59f))
## Other Changes
* style(api): refactoring for styling (ammarekbote)([5fdee63](https://github.com/lacework/terraform-provider-lacework/commit/5fdee63d718e3c25f58ea0d68a24b67c2f3e8c55))
* chore: update go-sdk dependency (#420) (Darren)([132d8fc](https://github.com/lacework/terraform-provider-lacework/commit/132d8fce5768e916345b6413d0c4ecadc894f8f6))
* ci: version bump to v1.1.1-dev (Lacework)([fd801d1](https://github.com/lacework/terraform-provider-lacework/commit/fd801d1ccc68d4b5d0f7b59290ab8240523cbb80))

## v1.1.0 (November 23, 2022)

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

## v1.0.1 (November 09, 2022)

## Bug Fixes
* fix: azure_cfg resource sening incorrect cloud account type (#403) (Darren)([1036ecd](https://github.com/lacework/terraform-provider-lacework/commit/1036ecd8ddc7f09b9ebefb7f8ba30d8fa8a02655))

## v1.0.0 (November 08, 2022)

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
* fix: remove unused func (Darren Murray)([dddeb21](https://github.com/lacework/terraform-provider-lacework/commit/dddeb2157f5d4e6aa997c464e2d91a9fb5a24f60))
* fix: remove unused funcs (Darren Murray)([b6ff8f8](https://github.com/lacework/terraform-provider-lacework/commit/b6ff8f89b88a9f698fa06d4cd2ade5fe7cdece6b))
## Documentation Updates
* docs: update deprecated cli cmd in docs (#396) (Darren)([1394822](https://github.com/lacework/terraform-provider-lacework/commit/139482264d9c1e5d516b90535c64e2f12a34ca65))
## Other Changes
* build(deps): bump golang.org/x/text from 0.3.7 to 0.4.0 (#388) (dependabot[bot])([83bccc0](https://github.com/lacework/terraform-provider-lacework/commit/83bccc07a080cbba2349432248a62550b8547601))
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#392) (dependabot[bot])([3fc1375](https://github.com/lacework/terraform-provider-lacework/commit/3fc1375cd2e119c1c19de719e41f53827f56a77f))
* build(deps): bump github.com/gruntwork-io/terratest (#386) (dependabot[bot])([eafb1b4](https://github.com/lacework/terraform-provider-lacework/commit/eafb1b4d1e903cb2eec670cacdc3ed9d525c935c))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#387) (dependabot[bot])([5dbbd34](https://github.com/lacework/terraform-provider-lacework/commit/5dbbd34f66e52f5838c9247e2635d8d92749ffea))
* build(deps): bump github.com/lacework/go-sdk from 0.43.0 to 0.44.1 (#393) (dependabot[bot])([43fb505](https://github.com/lacework/terraform-provider-lacework/commit/43fb505b7150a096acee33ce1d11fd9bef7f30e8))
* ci: version bump to v0.27.1-dev (Lacework)([2d8e4f8](https://github.com/lacework/terraform-provider-lacework/commit/2d8e4f882d35ee7d66d09c574e8269de3e75be6f))
* test: update integration test read to use v2 api (Darren Murray)([f5323f6](https://github.com/lacework/terraform-provider-lacework/commit/f5323f604b7be29d4b9ac38c33692a85561c26f0))
* test: move all tests to use APIv2 (#394) (Darren)([b88acee](https://github.com/lacework/terraform-provider-lacework/commit/b88acee715d181e0e20289a4106dd4ad56a7ba86))

## v0.27.0 (October 20, 2022)

## Features
* feat(resource): New lacework_integration_aws_org_agentless_scanning resource (#385) (Teddy Reed)([bf19236](https://github.com/lacework/terraform-provider-lacework/commit/bf19236f9bee889ae08c72a8924fce32e0f06104))
## Bug Fixes
* fix: suppressdiff for gcp_at private_key_id (#389) (Darren)([3f52c3b](https://github.com/lacework/terraform-provider-lacework/commit/3f52c3bf203cb20964354348e1057ef4ffeead04))
## Other Changes
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#368) (dependabot[bot])([9f48afe](https://github.com/lacework/terraform-provider-lacework/commit/9f48afe418ac56289d6df671ce67b85481a82003))
* ci: fix release pipeline by running make prepare (#390) (Darren)([94c4919](https://github.com/lacework/terraform-provider-lacework/commit/94c49193b8d94cd412f4af3a86396a85e17c537c))
* ci: version bump to v0.26.2-dev (Lacework)([7d46d13](https://github.com/lacework/terraform-provider-lacework/commit/7d46d13da7c57a87cf9aa6015b7c447d590b9143))

## v0.26.1 (September 26, 2022)

## Bug Fixes
* fix: suppressdiff for gcp_cfg private_key_id (#380) (Darren)([7188601](https://github.com/lacework/terraform-provider-lacework/commit/71886010cb7d3e68dbbdc2f8bb955908f64d8935))
* fix: fix gcr limit_by_labels field (#379) (Darren)([20a386c](https://github.com/lacework/terraform-provider-lacework/commit/20a386c82e8ac1f6bfd3d707540c96c9895de2e3))
## Other Changes
* ci: version bump to v0.26.1-dev (Lacework)([bd67f3b](https://github.com/lacework/terraform-provider-lacework/commit/bd67f3b039674c709f2c76e3e163008dce789390))

## v0.26.0 (September 21, 2022)

## Features
* feat(resource): New lacework_date_export_rule resource (#356) (Darren)([5d3b643](https://github.com/lacework/terraform-provider-lacework/commit/5d3b6431d3714e7649dc73ccec1eeaa86ee3eefb))
## Refactor
* refactor: migrate gcp_at to use v2 api (#369) (Darren)([41411b9](https://github.com/lacework/terraform-provider-lacework/commit/41411b900a8eb562a958d278a4f2633ef4df6086))
* refactor: migrate lacework_integration_gcp_cfg to use v2 api (#363) (Darren)([d21459d](https://github.com/lacework/terraform-provider-lacework/commit/d21459d0dd10aa317348be51d4b9c2cfdb7acd15))
* refactor: migrate lacework_integration_ecr to use v2 api (#364) (Darren)([88ee9bf](https://github.com/lacework/terraform-provider-lacework/commit/88ee9bf7a0cbe97391ef0aa5ea704c95c25a9d03))
* refactor: migrate lacework_integration_gcp_gcr to use v2 api (#365) (Darren)([3f2c581](https://github.com/lacework/terraform-provider-lacework/commit/3f2c5812d6ab7a4a7f3d1913132dcb86383c4fff))
* refactor: migrate lacework_integration_aws_cfg to use v2 api  (#361) (Darren)([f182648](https://github.com/lacework/terraform-provider-lacework/commit/f1826482cc217d82109b596a41f80a7d2f3734a3))
## Bug Fixes
* fix: fix go-sdk deps (#372) (Darren)([a1c9ff1](https://github.com/lacework/terraform-provider-lacework/commit/a1c9ff120f2a0bdfbc49cda81c6b6a250f8ca3a7))
## Other Changes
* ci: version bump to v0.25.1-dev (Lacework)([e2d2a1f](https://github.com/lacework/terraform-provider-lacework/commit/e2d2a1fced62deba901601a2083e0970220c8d97))

## v0.25.0 (August 16, 2022)

## Features
* feat(data-source): added lacework_user_profile (#357) (Alan Nix)([31380aa](https://github.com/lacework/terraform-provider-lacework/commit/31380aa3c1887c3dbc37f9cb8d249a4a115fb345))
## Other Changes
* ci: version bump to v0.24.1-dev (Lacework)([cc441cd](https://github.com/lacework/terraform-provider-lacework/commit/cc441cd4a31908f0d06e71170a765e3db6938e9a))

## v0.24.0 (August 09, 2022)

## Features
* feat: add credentials to aws_agentless_scanning resource (#354) (Darren)([c6426dd](https://github.com/lacework/terraform-provider-lacework/commit/c6426dd4525bd0f978cca7e7304983ce87efebcf))
## Other Changes
* ci: version bump to v0.23.1-dev (Lacework)([b288dfc](https://github.com/lacework/terraform-provider-lacework/commit/b288dfc389bd9374bb5f3e843a1f0aae927fae5c))

## v0.23.0 (August 08, 2022)

## Features
* feat(resource): New lacework_integration_aws_agentless_scanning resource   (#345) (Darren)([e4ee82f](https://github.com/lacework/terraform-provider-lacework/commit/e4ee82ff571f2ba373a2f8702b7eac2c5e5d6803))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#341) (dependabot[bot])([5cd8e26](https://github.com/lacework/terraform-provider-lacework/commit/5cd8e265a62be2fae8c3089f8ac5743078c65c82))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#350) (dependabot[bot])([455951e](https://github.com/lacework/terraform-provider-lacework/commit/455951e8877189f3ecdb951d0fa8be703527aa74))
* ci: version bump to v0.22.3-dev (Lacework)([6fc500e](https://github.com/lacework/terraform-provider-lacework/commit/6fc500ea0d162a5d9d35c3fa06b36905dbeac2f2))

## v0.22.2 (August 02, 2022)

## Bug Fixes
* fix: fixable parameter for vulnerability_exception resources (#342) (Darren)([4b5bb7a](https://github.com/lacework/terraform-provider-lacework/commit/4b5bb7a6751dccbda3532609b1c2f7cfdac9fec1))
## Other Changes
* build(deps): bump github.com/stretchr/testify from 1.7.3 to 1.8.0 (#337) (dependabot[bot])([2229d06](https://github.com/lacework/terraform-provider-lacework/commit/2229d069c887086f4a9e5a4f23fef6ce8b14b586))
* ci: version bump to v0.22.2-dev (Lacework)([a420033](https://github.com/lacework/terraform-provider-lacework/commit/a42003381d099ead77b5f4d1f652765f9ba5b8eb))

## v0.22.1 (July 22, 2022)

## Documentation Updates
* docs: adds api_token to argument reference (#348) (Salim Afiune)([2f5b07a](https://github.com/lacework/terraform-provider-lacework/commit/2f5b07a36cf0c41b01f8b390f09ebacc72ec9e32))
## Other Changes
* ci: version bump to v0.22.1-dev (Lacework)([f0e1b40](https://github.com/lacework/terraform-provider-lacework/commit/f0e1b40b6e4ed0ae16bba6941d305edc91875b47))
* test: Use unique API token for all integration tests (#347) (Salim Afiune)([0c85f89](https://github.com/lacework/terraform-provider-lacework/commit/0c85f892c5db68c76cc11fa0a3ef67be50399dd9))

## v0.22.0 (July 20, 2022)

## Features
* feat: support authentication via API token (#344) (Salim Afiune)([4550bce](https://github.com/lacework/terraform-provider-lacework/commit/4550bce99c97babc8cefef5caff62ed8018493a1))
* feat: new lacework_policy_exception resource (#326) (Darren)([01e87d3](https://github.com/lacework/terraform-provider-lacework/commit/01e87d37dd668910cdf89cc92c510e61ac0d73f1))
## Other Changes
* ci: version bump to v0.21.1-dev (Lacework)([7681f5d](https://github.com/lacework/terraform-provider-lacework/commit/7681f5d627112655270f526c66b299446d010eb7))

## v0.21.0 (June 30, 2022)

## Features
* feat(GcpGkeAudit): Add support for GcpGkeAudit apiv2 cloud account integration (#327) (Ross)([00d8d43](https://github.com/lacework/terraform-provider-lacework/commit/00d8d435585bda710d248d58782ae067673c7bc3))
## Bug Fixes
* fix(RAIN-33379): Fix documentation of policy type (#334) (Carsten Varming)([3fed26b](https://github.com/lacework/terraform-provider-lacework/commit/3fed26bbfd83c4d9fc6da3c356610c22eda650eb))
## Other Changes
* ci: version bump to v0.20.2-dev (Lacework)([6a64dfe](https://github.com/lacework/terraform-provider-lacework/commit/6a64dfe89da6f70b0ffdb0001ac198d78e0820c8))

## v0.20.1 (June 23, 2022)

## Bug Fixes
* fix: time parse error for expiry field (#332) (Darren)([8c691ef](https://github.com/lacework/terraform-provider-lacework/commit/8c691efe52d3d0f4e7ada25ce8c14ad2428bb911))
## Other Changes
* ci: version bump to v0.20.1-dev (Lacework)([5550234](https://github.com/lacework/terraform-provider-lacework/commit/555023422e27589e8448f8adb4828affa8b140bc))

## v0.20.0 (June 22, 2022)

## Features
* feat: add expiry field for Vulnerability Exceptions (#316) (Darren)([fb654c0](https://github.com/lacework/terraform-provider-lacework/commit/fb654c07a68e61b58660fa66f17b01894225c260))
## Bug Fixes
* fix: handle resource deleted outside of terraform scope (#331) (Darren)([b519479](https://github.com/lacework/terraform-provider-lacework/commit/b5194794f2d9a385ee467ca5901e2745d6be8628))
## Documentation Updates
* docs: update report_rule.html.markdown (#304) (Daniel Fitzgerald)([23b09f0](https://github.com/lacework/terraform-provider-lacework/commit/23b09f023f910546d44fed77a950d2851889c7a9))
## Other Changes
* build(deps): bump github.com/stretchr/testify from 1.7.1 to 1.7.3 (#330) (dependabot[bot])([4f00228](https://github.com/lacework/terraform-provider-lacework/commit/4f002285d489452334b294c1a659e0d62cf9eafe))
* build(deps): bump github.com/gruntwork-io/terratest (#329) (dependabot[bot])([a612993](https://github.com/lacework/terraform-provider-lacework/commit/a6129930cac59533211d7cdd1dd6625d0dc98e75))
* ci: fix deps (Salim Afiune Maya)([0cea8e3](https://github.com/lacework/terraform-provider-lacework/commit/0cea8e312e4df051fb0ac3b1a9544c6a9b5c4c3f))
* ci: version bump to v0.19.2-dev (Lacework)([2fa4909](https://github.com/lacework/terraform-provider-lacework/commit/2fa490936071012b243ebe8c88580ed346ee72e7))
* test(policy): fix policy type in integration testing (#328) (hazedav)([59f23e8](https://github.com/lacework/terraform-provider-lacework/commit/59f23e8563d6b6757fdc9ee8044a7eab7f89ef5a))

## v0.19.1 (June 07, 2022)

## Bug Fixes
* fix: allow resource scope to be optional (#320) (Salim Afiune)([501d0e7](https://github.com/lacework/terraform-provider-lacework/commit/501d0e75ddc543e64f54c7733acd7ada4c07d233))
## Other Changes
* ci: version bump to v0.19.1-dev (Lacework)([3e8a1d1](https://github.com/lacework/terraform-provider-lacework/commit/3e8a1d14becb539f1c40d9e12d22cb59764384c9))
* test: fix policy test (#312) (Darren)([d15112e](https://github.com/lacework/terraform-provider-lacework/commit/d15112e3832c61d7a0be8dcb050b2e98d03e6149))

## v0.19.0 (May 18, 2022)

## Features
* feat: add tags field to policy resource (#309) (Darren)([d353817](https://github.com/lacework/terraform-provider-lacework/commit/d353817f2b7bb2838e9bb2e9dc688ab72a1903e9))
## Other Changes
* ci: version bump to v0.18.1-dev (Lacework)([c434294](https://github.com/lacework/terraform-provider-lacework/commit/c43429431ec462dfc9b121f66a701bedff2f34c3))

## v0.18.0 (May 06, 2022)

## Features
* feat(resource): New lacework_alert_profile resource (#301) (Darren)([7812537](https://github.com/lacework/terraform-provider-lacework/commit/7812537d240b1016fe8d329fbb470a2b36433d54))
## Other Changes
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#296) (dependabot[bot])([cc8f492](https://github.com/lacework/terraform-provider-lacework/commit/cc8f492b8f76be6b1fb316fcebaeed2025e1c76c))
* ci: version bump to v0.17.1-dev (Lacework)([4769702](https://github.com/lacework/terraform-provider-lacework/commit/476970263e767a50f0d10e383ec4016c13bb9610))

## v0.17.0 (April 05, 2022)

## Features
* feat(AwsEksAudit): Add support for AwsEksAudit apiv2 cloud account integration (#292) (Ross)([4db5abe](https://github.com/lacework/terraform-provider-lacework/commit/4db5abeed889ac73df1066da27c9eedf50447da2))
* feat: add validation to ecr 'limit_num_imgs' (#283) (Darren)([a437790](https://github.com/lacework/terraform-provider-lacework/commit/a437790cf55d93ada23960066d1ed2197ee5fcac))
## Bug Fixes
* fix(docs): Fix aws eks audit log docs (#297) (Ross)([f9ab5cd](https://github.com/lacework/terraform-provider-lacework/commit/f9ab5cd7386c76146122ac15b4932898dad472f6))
* fix: migrate agent access token data source to use v2 api (#288) (Darren)([71cd925](https://github.com/lacework/terraform-provider-lacework/commit/71cd92593098f9b8c2cf3b79d27d25b302c1f15a))
## Documentation Updates
* docs: fix TF registry documentation tree (Salim Afiune Maya)([39cf2d7](https://github.com/lacework/terraform-provider-lacework/commit/39cf2d78d5bc51557de7c5a3054ba276f4e2a723))
* docs: fix typo (#282) (Darren)([deea6dd](https://github.com/lacework/terraform-provider-lacework/commit/deea6dd90d47c82765a0f40320be322783c071ae))
## Other Changes
* chore: update deps (#294) (Darren)([b6f3297](https://github.com/lacework/terraform-provider-lacework/commit/b6f32971e44d74b361c5688077054b6fdbd9c907))
* chore: update go version to 1.18 (#287) (Darren)([4618e3b](https://github.com/lacework/terraform-provider-lacework/commit/4618e3b02f97b15627ddb281c36f5293e0797528))
* chore: deprecate CLI Wiki (Salim Afiune Maya)([156ccaa](https://github.com/lacework/terraform-provider-lacework/commit/156ccaa48fc7fca13f127361934f79dc8eb4d6d7))
* chore: comment out error check in rollback test (#284) (Salim Afiune)([eaed849](https://github.com/lacework/terraform-provider-lacework/commit/eaed849420e5d71664572a5e862c002952b71fd1))
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#290) (dependabot[bot])([25059b5](https://github.com/lacework/terraform-provider-lacework/commit/25059b579a33ac0db74a714f843cd7b28f271f80))
* build(deps): bump github.com/gruntwork-io/terratest (#285) (dependabot[bot])([1e8ca59](https://github.com/lacework/terraform-provider-lacework/commit/1e8ca590f6b96691458c391b6b15986be0af4a6d))
* ci: add make cmd for output go tests in junit format (#293) (Darren)([eec76b2](https://github.com/lacework/terraform-provider-lacework/commit/eec76b2b3ef4d0c66baa469d6288c6daec075da8))
* ci: version bump to v0.16.1-dev (Lacework)([3f7eca5](https://github.com/lacework/terraform-provider-lacework/commit/3f7eca57944a51c9f4805fc40a52f06fae92ac19))

## v0.16.0 (March 02, 2022)

## Features
* feat: new lacework_query resource (#266) (Darren)([37377bb](https://github.com/lacework/terraform-provider-lacework/commit/37377bb8d1fb994da331d7bed8571fdce8c4ec7d))
* feat: new lacework_policy resource (#267) (Darren)([ab3efd8](https://github.com/lacework/terraform-provider-lacework/commit/ab3efd8d5cb973dde3036ce2105e0338e86b6ba0))
* feat: add registry notifications to docker v2 resource (#265) (Darren)([1650b9f](https://github.com/lacework/terraform-provider-lacework/commit/1650b9f801be012d89be05d424285cc8d1f1307c))
## Refactor
* refactor: encourage the use of anonymous queries (#280) (Salim Afiune)([b7c2d4c](https://github.com/lacework/terraform-provider-lacework/commit/b7c2d4ce9637606154293bf34101a78c078c4dc7))
* refactor: remove evaluator_id from policy and query resources (#278) (Salim Afiune)([953300a](https://github.com/lacework/terraform-provider-lacework/commit/953300a340879e3648c31180cba9e372047cc091))
## Bug Fixes
* fix: prevent updating evaluator_id on existing queries (#276) (Darren)([85dee76](https://github.com/lacework/terraform-provider-lacework/commit/85dee768fad25e73e0871d0eb80bece916cce134))
## Documentation Updates
* docs: update query and policy resources (#275) (Darren)([b85d686](https://github.com/lacework/terraform-provider-lacework/commit/b85d686d7c35db25ecb03a8422c8397587deef09))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#277) (dependabot[bot])([1444587](https://github.com/lacework/terraform-provider-lacework/commit/1444587d0ae8b16130015270ebd0b78ad614e3e4))
* build(deps): bump github.com/gruntwork-io/terratest (#271) (dependabot[bot])([edb3af8](https://github.com/lacework/terraform-provider-lacework/commit/edb3af87f74f21d32183eac24c82166a22ac8c18))
* build(deps): bump github.com/gruntwork-io/terratest (#262) (dependabot[bot])([30df648](https://github.com/lacework/terraform-provider-lacework/commit/30df648c42cadffbeaa777e285a9ea39034cbe85))
* ci: version bump to v0.15.1-dev (Lacework)([07ced51](https://github.com/lacework/terraform-provider-lacework/commit/07ced513d675b1493593d6b2df71b0dbfe1c400e))
* test(team_members): use lwdomain to get correct account name (#272) (Darren)([e120def](https://github.com/lacework/terraform-provider-lacework/commit/e120def7f14f533b867fe904610c906ec5168d1e))
* test: skip org test for lw account resource group (#261) (Salim Afiune)([061ada9](https://github.com/lacework/terraform-provider-lacework/commit/061ada9f9c0ecf4821ab1dd9dc7633a9a88db06f))

## v0.15.0 (January 28, 2022)

## Features
* feat(resource): New lacework_vulnerability_exception_container (#253) (Darren)([6903c79](https://github.com/lacework/terraform-provider-lacework/commit/6903c79e6270f5b1dc1e5a06499f4804a3f49ba3))
* feat(resource): New lacework_vulnerability_exception_host (#248) (Darren)([afa657a](https://github.com/lacework/terraform-provider-lacework/commit/afa657aaea7bc17e75d2da61920bf39935fda616))
## Refactor
* refactor: remove expiration time (#254) (Darren)([6bbc742](https://github.com/lacework/terraform-provider-lacework/commit/6bbc7427cee5946a2288369d129250d5f9088abf))
## Documentation Updates
* docs: fix PR template document (#255) (Darren)([bfad072](https://github.com/lacework/terraform-provider-lacework/commit/bfad0729b364ec49fae8b3dc916f4998e52bae24))
* docs: update 'aws' s3/cloudwatch to 'amazon' in documentation (#256) (Darren)([1a3bece](https://github.com/lacework/terraform-provider-lacework/commit/1a3bece85234560450e382ba4f2a5e01492959ef))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#246) (dependabot[bot])([94ff17f](https://github.com/lacework/terraform-provider-lacework/commit/94ff17f1292c4e4348a70340985b118d9f4888da))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#251) (dependabot[bot])([8a9dbca](https://github.com/lacework/terraform-provider-lacework/commit/8a9dbcacfe38e26395d915f2f75220e7e8ffda92))
* ci: version bump to v0.14.1-dev (Lacework)([9e9b210](https://github.com/lacework/terraform-provider-lacework/commit/9e9b2109371bb95c508823ec907910c596cf918d))

## v0.14.0 (December 17, 2021)

## Features
* feat: new lacework_team_member resource (#245) (vatasha)([cb4e69a](https://github.com/lacework/terraform-provider-lacework/commit/cb4e69a4ceb3138d9874589b934d502ecfd3f776))
## Refactor
* refactor: enable non-OS Package Support by default (#249) (robewedd)([7c85cbb](https://github.com/lacework/terraform-provider-lacework/commit/7c85cbbd029ac7e88b3cc6184b9a60d6b252372d))
## Bug Fixes
* fix: remove apiIntgKey from resourceLaceworkAlertChannelPagerDutyRead (#247) (Darren)([ad170ef](https://github.com/lacework/terraform-provider-lacework/commit/ad170efcd17995b8c9957631ae947913dc259ffb))
## Other Changes
* ci: version bump to v0.13.1-dev (Salim Afiune Maya)([e5800e0](https://github.com/lacework/terraform-provider-lacework/commit/e5800e0bac24598319d95945321d7e424a7adfff))
* test: Use v2 endpoints in tests (#244) (Darren)([ba89973](https://github.com/lacework/terraform-provider-lacework/commit/ba899737c17d37ce180b24799836caf21c8c362d))

## v0.13.0 (December 06, 2021)

## Features
* feat(resource): New lacework_report_rule (#237) (Darren)([c2928b6](https://github.com/lacework/terraform-provider-lacework/commit/c2928b612dc7e1db2c676bbd9abcf4f2dce7348e))
## Documentation Updates
* docs: mention new GAR and GCR Modules (#229) (Salim Afiune)([30d1c0a](https://github.com/lacework/terraform-provider-lacework/commit/30d1c0a8b5e305ccf27ffc8e31f9b541ad78033a))
* docs: Add environment variables for Windows (#228) (Salim Afiune)([bd4ee0c](https://github.com/lacework/terraform-provider-lacework/commit/bd4ee0cdd86da76582fd870aa3fc2346c22b30d4))
## Other Changes
* chore: run make go-vendor (#241) (Darren)([3dd3649](https://github.com/lacework/terraform-provider-lacework/commit/3dd364912e64c9fc44242ccbf97030ff382535a4))
* ci: run vendor commands after dep update (#227) (Salim Afiune)([8b18fb8](https://github.com/lacework/terraform-provider-lacework/commit/8b18fb8837345234bd074579592c5c64741fe771))
* ci: version bump to v0.12.3-dev (Lacework)([528641b](https://github.com/lacework/terraform-provider-lacework/commit/528641ba9dbd7dd23a1c506bb427c829ec0deff1))
* test: disable failing alert rule tests (#240) (Darren)([4247ec9](https://github.com/lacework/terraform-provider-lacework/commit/4247ec97aead69d5ac4fe5e2f2ece61c1dfa6f67))
* test: create unique names for resource groups integration tests (#239) (Darren)([ed382b1](https://github.com/lacework/terraform-provider-lacework/commit/ed382b1e3ad1a394c7349e043a5c694012540657))
* test: fix alert-rule tests (#236) (Darren)([c92b8da](https://github.com/lacework/terraform-provider-lacework/commit/c92b8da3cf0b0b5f80a8c46859f76ebd71c65394))
* test: fix s3 alert channel integration test (#224) (vatasha)([5f2cad5](https://github.com/lacework/terraform-provider-lacework/commit/5f2cad56af47926305085846c8c5fabf0b5b7405))

## v0.12.2 (November 10, 2021)

## Bug Fixes
* fix: Deprecate channels attribute in favor of alert_channels (#225) (Darren)([ba30502](https://github.com/lacework/terraform-provider-lacework/commit/ba305020a31cf486a4baa1bf7079b65b446822b4))
* fix: Suppress diff on sensitive values (#221) (Darren)([4a569cd](https://github.com/lacework/terraform-provider-lacework/commit/4a569cd485534d54ce4938e283ad267e9252fe4e))
## Other Changes
* style: fix linting and run it in CI(#222) (Darren)([725649a](https://github.com/lacework/terraform-provider-lacework/commit/725649ac626f52bfc1036475a37d8616c864808f))
* ci: version bump to v0.12.2-dev (Lacework)([cc7609b](https://github.com/lacework/terraform-provider-lacework/commit/cc7609bcc7aece7f4e28f4192c54206ba72132a8))

## v0.12.1 (October 29, 2021)

## Documentation Updates
* docs: Fix alert rule doc summary (#218) (Darren)([e662841](https://github.com/lacework/terraform-provider-lacework/commit/e662841024f340ae4b8d3ee4e607dcd8ad104834))
* docs: Fix alert rule page name (#216) (Darren)([d6f1c3b](https://github.com/lacework/terraform-provider-lacework/commit/d6f1c3bc03e425a6a1565758e73df92667a8c948))
## Other Changes
* ci: improve unit and integration test output (#217) (Salim Afiune)([518b106](https://github.com/lacework/terraform-provider-lacework/commit/518b106f558ebf725e8a5297889d07ae0cd98cab))
* ci: version bump to v0.12.1-dev (Lacework)([4eefa67](https://github.com/lacework/terraform-provider-lacework/commit/4eefa67191d958448e31016bc7bacbb106662876))

## v0.12.0 (October 29, 2021)

## Features
* feat(resource): New lacework_alert_rule (#211) (Darren)([c62e1e1](https://github.com/lacework/terraform-provider-lacework/commit/c62e1e17d684222333ea9c986cb71f7f1ee6c586))
## Bug Fixes
* fix: Migrate Jira Cloud and Jira Server to API v2 (#212) (vatasha)([2382376](https://github.com/lacework/terraform-provider-lacework/commit/2382376940504453576ab1298f21c3e45021289c))
* fix: Migrate IbmQRadar alert channel(v2) (#208) (Darren)([2699390](https://github.com/lacework/terraform-provider-lacework/commit/2699390cfd2eb11b7e638edc1d69ac7982698f7f))
* fix: Migrate PagerDuty alert channel(v2) (#209) (Darren)([29057fc](https://github.com/lacework/terraform-provider-lacework/commit/29057fc71c9eb90af223c6fea740978c5feed318))
* fix: Migrate NewRelic Insights alert channel(v2) (#210) (Darren)([553d395](https://github.com/lacework/terraform-provider-lacework/commit/553d395863f8171d6da23add44993984da305e02))
* fix: Migrate ServiceNow alert channel(v2) (#206) (Darren)([09a413e](https://github.com/lacework/terraform-provider-lacework/commit/09a413ea6757f07662d26af8262d2290858df28e))
* fix: Migrate Splunk alert channel(v2) (#205) (Darren)([7cc2d5e](https://github.com/lacework/terraform-provider-lacework/commit/7cc2d5e59d17bc1381ed8c227b48ce776ab4c474))
* fix: Migrate GCP Pub Sub alert channel to API v2 (#207) (vatasha)([cd3da9a](https://github.com/lacework/terraform-provider-lacework/commit/cd3da9a31c58f247d35382288af50bb35d86ab6c))
## Documentation Updates
* docs: fix webhook alert channel example (#204) (Salim Afiune)([fba343b](https://github.com/lacework/terraform-provider-lacework/commit/fba343bb8d8d62f2b2e6abf4655bda69855c8d83))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#196) (dependabot[bot])([6b94df3](https://github.com/lacework/terraform-provider-lacework/commit/6b94df3d8bdb9eb79a41d10f85ef4911c14a4e71))
* ci: fix release pipeline (#214) (Salim Afiune)([bde019d](https://github.com/lacework/terraform-provider-lacework/commit/bde019d4dfa7cb98717d7179487080cb4a576cee))
* ci: version bump to v0.11.3-dev (Lacework)([7551e99](https://github.com/lacework/terraform-provider-lacework/commit/7551e99525ca9395ebedb53dc5e1a6a67625ba02))
* test: fix tests for alert_rule resource (#213) (Salim Afiune)([0f7aa5b](https://github.com/lacework/terraform-provider-lacework/commit/0f7aa5b513e95081e9e875fc5871440462e8ab1c))

## v0.11.2 (October 12, 2021)

## Bug Fixes
* fix: migrate Webhook alert channel(v2) (#188) (Darren)([a27565c](https://github.com/lacework/terraform-provider-lacework/commit/a27565c7a94a6faac34f1389b56d79a34dbd3a29))
* fix: migrate VictorOps alert channel(v2) (#192) (Darren)([517e587](https://github.com/lacework/terraform-provider-lacework/commit/517e58742b1371fe6406c24b66431793ab107b99))
* fix: migrate Microsoft Teams Alert Channels to API v2 (Salim Afiune Maya)([36a78e3](https://github.com/lacework/terraform-provider-lacework/commit/36a78e3473aa3f92a1acc20bc7c00347f87a361b))
* fix: migrate Cisco webex alert channel to API v2 (#193) (vatasha)([66e8251](https://github.com/lacework/terraform-provider-lacework/commit/66e8251d8fc252ee7313703bb4a6fde030cca912))
* fix: remove alert channel properly if test fails (#199) (Salim Afiune)([9e13d58](https://github.com/lacework/terraform-provider-lacework/commit/9e13d58b3af1bf1e5fc66575747c86b58142aaca))
## Other Changes
* chore(deps): update lacework/go-sdk from main (Salim Afiune Maya)([5d277f9](https://github.com/lacework/terraform-provider-lacework/commit/5d277f9c87b727c865feee54c6eedd90f6c99c46))
* ci: fix release pipeline (#202) (Salim Afiune)([a37168a](https://github.com/lacework/terraform-provider-lacework/commit/a37168a1949c1c0c912ec58696c3d3ad23dd6809))
* ci: version bump to v0.11.2-dev (Lacework)([7de4f58](https://github.com/lacework/terraform-provider-lacework/commit/7de4f580e8a43900b3b6f5afc2fe51b49e2d6a5a))
* test: fix tf output values (#198) (Darren)([a4af709](https://github.com/lacework/terraform-provider-lacework/commit/a4af709c3eabe6471bc3f9e0258a30201231ad82))

## v0.11.1 (October 11, 2021)

## Bug Fixes
* fix: increase timeout from 60 to 125 seconds (#194) (Salim Afiune)([ff3819e](https://github.com/lacework/terraform-provider-lacework/commit/ff3819e782167d3a7d987b71018124c2ddedf6ec))
* fix: Migrate Datadog alert channel to API v2 (#189) (vatasha)([fd38f37](https://github.com/lacework/terraform-provider-lacework/commit/fd38f37ca5e3851a08c7d2bd522cdaf5e3d202c5))
* fix: migrate AWS Cloudwatch alert channel to API v2 (#186) (vatasha)([eeb55a7](https://github.com/lacework/terraform-provider-lacework/commit/eeb55a7000e214374abcb1bf930ee724ba588a84))
## Other Changes
* chore(deps): update lacework/go-sdk from main (Salim Afiune Maya)([37245e8](https://github.com/lacework/terraform-provider-lacework/commit/37245e8e5d0ccc44007f38e69335263887c45f92))
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#176) (dependabot[bot])([5800a22](https://github.com/lacework/terraform-provider-lacework/commit/5800a228802f076a4dd584e2022d3bf56a874cae))
* build(deps): bump github.com/gruntwork-io/terratest (#187) (dependabot[bot])([994389b](https://github.com/lacework/terraform-provider-lacework/commit/994389b0cbbdc1032f2322cf3cfb89e7a3b20444))
* ci: fix downgrading go packages (Salim Afiune Maya)([6eb22dc](https://github.com/lacework/terraform-provider-lacework/commit/6eb22dc1adaf1d2f31f1d2987cd3b0af43f14fed))
* ci: version bump to v0.11.1-dev (Lacework)([6d86c83](https://github.com/lacework/terraform-provider-lacework/commit/6d86c837086919b3d823993396647055dfd8f50b))

## v0.11.0 (October 01, 2021)

## Features
* feat: Add new field non_os_package_support to ECR resource (#175) (Darren)([f150be8](https://github.com/lacework/terraform-provider-lacework/commit/f150be8ad4163ace9acefa3244145aaed22f7c75))
* feat: Add new field non_os_package_support to GHCR resource (#178) (Darren)([e72b5eb](https://github.com/lacework/terraform-provider-lacework/commit/e72b5eb702385103a9384dc10ec37255a9d98168))
* feat: Add new field non_os_package_support to GAR resource (#179) (Darren)([27b6c70](https://github.com/lacework/terraform-provider-lacework/commit/27b6c70bdddf154f8171e5bf71d929bcea015ca4))
## Documentation Updates
* docs: update multiple aliases for sub-accounts (#172) (robewedd)([0763028](https://github.com/lacework/terraform-provider-lacework/commit/0763028b8b4d4a5df29ab339003ba8ba5a680d91))
## Other Changes
* ci: fix release pipeline by running make prepare (#183) (Salim Afiune)([fb4d3db](https://github.com/lacework/terraform-provider-lacework/commit/fb4d3db24cab8a94e8e8ccaf8086dd44bae3667d))
* ci: version bump to v0.10.1-dev (Lacework)([9645ee1](https://github.com/lacework/terraform-provider-lacework/commit/9645ee118106c2358658c955bc50f892cabb0153))
* test: fix integration tests (#180) (Darren)([1eb35e9](https://github.com/lacework/terraform-provider-lacework/commit/1eb35e9fee926dab0caa25e9a55bae95a3a47076))

## v0.10.0 (September 27, 2021)

## Features
* feat: New Lw Account Resource Group Terraform Resource (#171) (Darren)([8425534](https://github.com/lacework/terraform-provider-lacework/commit/84255348e6e68986ccdf773ce73219990032dac5))
* feat: New Container Resource Group Terraform Resource (#170) (Darren)([ffc35cf](https://github.com/lacework/terraform-provider-lacework/commit/ffc35cf73537b0968619831e873d92d0ca6df4da))
* feat: New Machine Resource Group Terraform Resource (#169) (Darren)([4433d9a](https://github.com/lacework/terraform-provider-lacework/commit/4433d9a10367c49b412a399d5358f96bb95b51f5))
## Refactor
* refactor: switch over to use APIv2 by default (#173) (Salim Afiune)([6ebd5fc](https://github.com/lacework/terraform-provider-lacework/commit/6ebd5fc0b3ebcc75fa00986633dfab66a35da1e4))
## Documentation Updates
* docs: Add contributing documentation (#165) (Darren)([55a4408](https://github.com/lacework/terraform-provider-lacework/commit/55a44080a0f2901500f5d791be357a7b163019b6))
## Other Changes
* ci: version bump to v0.9.3-dev (Lacework)([f619f25](https://github.com/lacework/terraform-provider-lacework/commit/f619f2572e3812b6db6257ffafebe8d4c7f4ec5f))

## v0.9.2 (September 23, 2021)

## Bug Fixes
* fix: Fix return empty array when casting limit_by fields (#167) (Darren)([329a00c](https://github.com/lacework/terraform-provider-lacework/commit/329a00cc21e3b2729e933df3fa271ad4a4b8fa29))
## Other Changes
* build(deps): bump github.com/gruntwork-io/terratest (#164) (dependabot[bot])([2d9be48](https://github.com/lacework/terraform-provider-lacework/commit/2d9be48831618f6cbba6a49be3c0e3f04e2f8eb6))
* ci: version bump to v0.9.2-dev (Lacework)([80e0cb8](https://github.com/lacework/terraform-provider-lacework/commit/80e0cb830d18801399381adf39d29d0893a8cc45))

## v0.9.1 (September 13, 2021)

## Documentation Updates
* docs: standardize mention of 'ids' (#162) (Salim Afiune)([8591935](https://github.com/lacework/terraform-provider-lacework/commit/8591935b8372dec0cdcd93575074b62356f48841))
* docs: fix html markdown indentation (#161) (Darren)([a8e735a](https://github.com/lacework/terraform-provider-lacework/commit/a8e735a741795744281f04611cfe69e7ab9e25cb))
## Other Changes
* ci: version bump to v0.9.1-dev (Lacework)([ec0b0b7](https://github.com/lacework/terraform-provider-lacework/commit/ec0b0b7cadd8be729ee167443da34de237f40c6f))

## v0.9.0 (September 10, 2021)

## Features
* feat(resource): New lacework_resource_group_azure (#158) (Darren)([6ab2f0c](https://github.com/lacework/terraform-provider-lacework/commit/6ab2f0c75b50e533f987bd31e51ce90ef052d503))
* feat(resource): New lacework_resource_group_gcp (#156) (Darren)([d88a7b9](https://github.com/lacework/terraform-provider-lacework/commit/d88a7b9671dc759c45e56ac2877694dfead7e278))
* feat: gracefully handle account config <ACCOUNT>.lacework.net (#157) (Salim Afiune)([cb32670](https://github.com/lacework/terraform-provider-lacework/commit/cb326706d479f5f405dd30c6414ed271e5f38bad))
* feat: Add Non-OS Package support for GCR, DockerV2, and DockerHub (#152) (Andre Elizondo)([96b4df8](https://github.com/lacework/terraform-provider-lacework/commit/96b4df89fcb52b1cffe6457e5db3df85e219a994))
* feat(resource): New lacework_resource_group_aws (#155) (Darren)([ca4eb3d](https://github.com/lacework/terraform-provider-lacework/commit/ca4eb3d3556b7dfad07ec749296cc67d76ac8be5))
* feat: new Github Container Registry (GHCR) resource (#143) (Darren)([7301f74](https://github.com/lacework/terraform-provider-lacework/commit/7301f7461d332d50184868d0ed6450b12d0a5bbb))
## Bug Fixes
* fix: enable parsing internal accounts (#159) (Salim Afiune)([bb72fb4](https://github.com/lacework/terraform-provider-lacework/commit/bb72fb410b7a4ed5ef1ae2851fcbd0bea17d825b))
## Other Changes
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 (#154) (dependabot[bot])([0e749a3](https://github.com/lacework/terraform-provider-lacework/commit/0e749a3e4b89c3f61ec5ea14df1649bade666e3a))
* build(deps): bump github.com/gruntwork-io/terratest (#151) (dependabot[bot])([3a8eaed](https://github.com/lacework/terraform-provider-lacework/commit/3a8eaed311a38b42e64d6b3050bd19bd6c605b62))
* ci: version bump to v0.8.2-dev (Lacework)([b998522](https://github.com/lacework/terraform-provider-lacework/commit/b998522ae2ed362f64ed53736999b8fd4a6f1187))

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
