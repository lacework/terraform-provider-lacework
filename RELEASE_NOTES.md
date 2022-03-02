# Release Notes
Another day, another release. These are the release notes for the version `v0.16.0`.

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
