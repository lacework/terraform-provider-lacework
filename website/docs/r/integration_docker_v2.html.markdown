---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_docker_v2"
description: |-
  Create and manage Docker V2 container registry integrations
---

# lacework\_integration\_docker\_v2

Use the Docker V2 Registry integration for private Docker V2 registries only.

~> **Note:** For Docker Hub, ECR, and GCR, use their corresponding container registry types.

The Docker V2 Registry integration functions differently than Lacework's other container registry
integrations. This integration performs on-demand image assessment via the Lacework API, while the other
integrations automatically assess images at regular intervals.

Supported Docker V2 registries:

* Azure Container Registry
* GitLab (On prem 12.8 and cloud)
* JFrog Artifactory (On prem 7.2.1 and cloud)
* JFrog Platform (On prem 7.2.1 and cloud)

~> **Note:** You must whitelist the Lacework outbound IPs to allow the vulnerability scanner to communicate with your private registries. See [Lacework Outbound IPs](https://support.lacework.com/hc/en-us/articles/360052140433)

## Set Up Image Assessments

The Lacework CLI makes it easy to request on-demand scans of new images designed for continuous
integration (CI) pipelines. You can find more information on integrating the Lacework CLI for
container vulnerability scanning in CI pipelines [here](https://support.lacework.com/hc/en-us/articles/360052476154-Integrate-Lacework-APIs-with-Continuous-Integration-CI-Pipelines).

-> **Note:** The Docker V2 Registry status displays `Integration Successful` only after its first assessment completes.

For more information visit the [documentation for the Lacework CLI](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#container-vulnerability-assessments).

## Example Usage

```hcl
resource "lacework_integration_docker_v2" "jfrog" {
  name = "My Docker V2 Registry"
  registry_domain = "my-dockerv2.jfrog.io"
  username = "my-user"
  password = "a-secret-password"
  ssl = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Docker V2 Registry integration name.
* `registry_domain` - (Required) The registry domain. Allowed formats are `YourIP:YourPort` or `YourDomain:YourPort`.
* `username` - (Required) The user that has at permissions to pull from the container registry the images to be assessed.
* `password` - (Required) The password for the specified user.
* `ssl` - (Optional) Enable or disable SSL communication. Defaults to `false`.
* `limit_by_tag` - (Optional, **Deprecated**) An image tag to limit the assessment of images with matching tag. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. **This attribute will be removed in version 1.0 of the Lacework provider.**
* `limit_by_label` - (Optional, **Deprecated**) An image label to limit the assessment of images with matching label. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. **This attribute will be removed in version 1.0 of the Lacework provider.**
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework Docker V2 container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_docker_v2.jfrog EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
