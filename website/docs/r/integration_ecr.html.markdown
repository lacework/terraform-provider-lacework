---
subcategory: "Container Registry Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_ecr"
description: |-
  Create and manage ECR integrations
---

# lacework\_integration\_ecr

Use this resource to integrate an Amazon Container Registry (ECR) with Lacework to assess, identify,
and report vulnerabilities found in the operating system software packages in a Docker container
image.

~> **Note:** Assessing a retagged ECR image is not supported because ECR does not consider it a new image and does not create a new entry. To assess a retagged image, use on-demand assessment through the Lacework CLI. For more information, see the [container vulnerability section in the Lacework CLI documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#container-vulnerability-assessments).

## Example Usage

```hcl
resource "lacework_integration_ecr" "example" {
  name              = "ERC Example"
  registry_domain   = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  access_key_id     = "AWS123abcAccessKeyID"
  secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The ECR integration name.
* `registry_domain` - (Required) The Amazon Container Registry (ECR) domain in the format `YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com`, where `YourAWSAcount` is the AWS account number for the AWS IAM user that has a role with permissions to access the ECR and `YourRegion` is your AWS region such as `us-west-2`.
* `access_key_id` - (Required) The AWS access key ID for an AWS IAM user that has a role with permissions to access the Amazon Container Registry (ECR).
* `secret_access_key` - (Required) The AWS secret key for the specified AWS access key.
* `limit_by_tag` - (Optional) An image tag to limit the assessment of images with matching tag. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`.
* `limit_by_label` - (Optional) An image label to limit the assessment of images with matching label. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`.
* `limit_by_repos` - (Optional) A comma-separated list of repositories to assess. (without spaces recommended)
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

## Import

A Lacework ECR integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_ecr.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).


