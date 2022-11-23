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

~> **Note:** You must whitelist the Lacework outbound IPs to allow the vulnerability scanner to communicate with your private registries. See [Lacework Outbound IPs](https://docs.lacework.com/lacework-outbound-ips)

## Example Usage

```hcl
resource "lacework_integration_docker_v2" "jfrog" {
  name            = "My Docker V2 Registry"
  registry_domain = "my-dockerv2.jfrog.io"
  username        = "my-user"
  password        = "a-secret-password"
  ssl             = true
}
```

-> **Note:** The Docker V2 Registry status displays `Integration Successful` only after its first assessment completes.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Docker V2 Registry integration name.
* `registry_domain` - (Required) The registry domain. Allowed formats are `YourIP:YourPort` or `YourDomain:YourPort`.
* `username` - (Required) The user that has at permissions to pull from the container registry the images to be assessed.
* `password` - (Required) The password for the specified user.
* `ssl` - (Optional) Enable or disable SSL communication. Defaults to `false`.
* `notifications` - (Optional) Subscribe to registry notifications. Defaults to `false`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `non_os_package_support` - (Optional) Enable [program language scanning](https://docs.lacework.com/container-image-support#language-libraries-support). Defaults to `true`.
* `limit_by_tags` - (Optional) A list of image tags to limit the assessment of images with matching tags. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_label` - (Optional) A list of key/value labels to limit the assessment of images. If you specify `limit_by_tags` and `limit_by_label` limits, they function as an `AND`.

The `limit_by_label` block can be defined multiple times to define multiple label limits, it supports:
* `key` - (Required) The key of the label.
* `value` - (Required) The value of the label.

## Import

A Lacework Docker V2 container registry integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_docker_v2.jfrog EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retrieve the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework container-registry list`. To install this tool follow
	[this documentation](https://docs.lacework.com/cli/).
