# Release Notes
Another day, another release. These are the release notes for the version `v2.5.0`.

## Features
* feat: add the ability to debug the provider (#747) (Pengyuan Zhao)([fe2ebce4](https://github.com/lacework/terraform-provider-lacework/commit/fe2ebce4b501c57515bff21606a27fa32ce86e8a))
* feat(RAIN-96864)!: remove decommissioned lacework_integration_gcp_at resource (#729) (Manan Bhatia)([3f082e1b](https://github.com/lacework/terraform-provider-lacework/commit/3f082e1ba632ddfb983536c5c4d924cdffa036aa))
## Bug Fixes
* fix(test): match new StringInSlice error format from plugin-sdk v2.40.1 (#745) (Pengyuan Zhao)([94b35c29](https://github.com/lacework/terraform-provider-lacework/commit/94b35c29a7b0b8effdbb1a03517a3cc9adf0cca9))
* fix(AWLS2-1027): DSPM integration_level property must be marked 'computed' (#728) (Joseph Wilder)([368ffeb8](https://github.com/lacework/terraform-provider-lacework/commit/368ffeb8965c53f6f99f5d76e8c260bc0a4e5fc1))
## Other Changes
* build(deps): bump go-sdk to v2.15.1 and testify to v1.11.1 (#741) (Pengyuan Zhao)([2fbffaef](https://github.com/lacework/terraform-provider-lacework/commit/2fbffaef2816f96e4783426f4e4519e6206e8280))
* build(deps): bump terraform-plugin-sdk from v2.27.0 to v2.40.1 (#740) (Pengyuan Zhao)([45507d25](https://github.com/lacework/terraform-provider-lacework/commit/45507d2558979264b218bbae0daf136fdb2836b3))
* build(deps): bump x/crypto, x/net, x/text, grpc for security advisories (#739) (Pengyuan Zhao)([cf716e35](https://github.com/lacework/terraform-provider-lacework/commit/cf716e3545f008356749ed4688654dc442eb85ff))
* build(deps): bump golang.org/x/crypto, golang.org/x/text, ulikunitz/xz (#730) (Lokesh Vadlamudi)([f867f653](https://github.com/lacework/terraform-provider-lacework/commit/f867f653303741cc7effcba2b4ebc1cfe49e1594))
* ci: bump actions/setup-go to v6 so the toolchain directive is honoured (#746) (Pengyuan Zhao)([89fb7544](https://github.com/lacework/terraform-provider-lacework/commit/89fb7544c46b36e54c99cd7615533e36a813f350))
* ci: group Dependabot updates for golang.org/x and terraform-plugin-* (#742) (Pengyuan Zhao)([0262aac9](https://github.com/lacework/terraform-provider-lacework/commit/0262aac9970cb0074b0b8958f233b22abf2234e9))
* ci: install Go toolchain from go.mod instead of pinning 1.21.x (#738) (Pengyuan Zhao)([259418b5](https://github.com/lacework/terraform-provider-lacework/commit/259418b5d7b8de8dbef906d2cdded38b59f81b40))
* ci: version bump to v2.4.1-dev (Lacework)([6df073d1](https://github.com/lacework/terraform-provider-lacework/commit/6df073d15a3b2797d4b18cd5e65a761a4cb724e5))
