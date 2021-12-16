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

~> **Note:** Assessing a retagged ECR image is not supported because ECR does not consider it a new
image and does not create a new entry. To assess a retagged image, use on-demand assessment through
the Lacework CLI. For more information, see the [container vulnerability section in the Lacework CLI
documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#container-vulnerability-assessments).

This resource has two authentication methods:

* AWS Access Key-Based Authentication
* AWS IAM Role-Based Authentication

!> **Warning:** It is possible to switch authentication methods but the resource
will be destroyed and recreated. This will generate a new `INT_GUID`.

For more information, see [Integrate Amazon Container Registry documentation](https://docs.lacework.com/integrate-amazon-container-registry)

## Example Usage

### Authentication via AWS Access Key
```hcl
resource "lacework_integration_ecr" "access_key" {
  name            = "ECR using Access Keys"
  non_os_package_support = true
  registry_domain = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  credentials {
    access_key_id     = "AWS123abcAccessKeyID"
    secret_access_key = "AWS123abc123abcSecretAccessKey0000000000"
  }
}
```

### Authentication via AWS IAM Role
```hcl
resource "lacework_integration_ecr" "iam_role" {
  name            = "ECR using IAM Role"
  registry_domain = "YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com"
  credentials {
    role_arn    = "arn:aws:iam::1234567890:role/lacework_iam_example_role"
    external_id = "12345"
  }
}
```

## ECR Module: Create an ECR Integration with an IAM Role

This example shows how to leverage the [Lacework ECR Terraform Module](https://registry.terraform.io/modules/lacework/ecr/aws/latest)
to automatically create a new IAM role and use it to create an ECR integration:

```hcl
provider "lacework" {}

provider "aws" {
  region = "us-west-2"
}

module "lacework_ecr" {
  source  = "lacework/ecr/aws"
  version = "~> 0.1"
  non_os_package_support = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The ECR integration name.
* `registry_domain` - (Required) The Amazon Container Registry (ECR) domain in the format `YourAWSAccount.dkr.ecr.YourRegion.amazonaws.com`, where `YourAWSAcount` is the AWS account number for the AWS IAM user that has a role with permissions to access the ECR and `YourRegion` is your AWS region such as `us-west-2`.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `limit_by_tag` - (Optional, **Deprecated**) An image tag to limit the assessment of images with matching tag. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. This attribute will be replaced by a new attribute `limit_by_tags` in version 1.0 of the Lacework provider.
* `limit_by_label` - (Optional, **Deprecated**) An image label to limit the assessment of images with matching label. If you specify `limit_by_tag` and `limit_by_label` limits, they function as an `AND`. Supported field input are `mytext*mytext`, `mytext`, `mytext*`, or `mytext`. Only one `*` wildcard is supported. Defaults to `*`. This attribute will be replaced by a new attribute `limit_by_labels` in version 1.0 of the Lacework provider.
* `limit_by_repos` - (Optional, **Deprecated**) A comma-separated list of repositories to assess. (without spaces recommended) This attribute will be replaced by a new attribute `limit_by_repositories` in version 1.0 of the Lacework provider.
* `limit_num_imgs` - (Optional) The maximum number of newest container images to assess per repository. Must be one of `5`, `10`, or `15`. Defaults to `5`.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.
* `limit_by_tags` - (Optional) A list of image tags to limit the assessment of images with matching tags. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_labels` - (Optional) A key based map of labels to limit the assessment of images with matching `key:value` labels. If you specify `limit_by_tags` and `limit_by_labels` limits, they function as an `AND`.
* `limit_by_repositories` - (Optional) A list of repositories to assess.
* `non_os_package_support` - (Optional) Enable [program language scanning](https://docs.lacework.com/container-image-support#language-libraries-support). Defaults to `true`.

### Credentials

`credentials` supports the combination of the following arguments.

**For AWS IAM Role-Based Authentication, only both of these arguments are required:**
* `role_arn` - The ARN of the IAM role with permissions to access the Amazon Container Registry (ECR).
* `external_id` - The external ID for the IAM role.

**For AWS Access Key-Based Authentication, only both of these arguments are required:**
* `access_key_id` - The AWS access key ID for an AWS IAM user that has a role with permissions to access the Amazon Container Registry (ECR).
* `secret_access_key` - The AWS secret key for the specified AWS access key.

## Import

A Lacework ECR integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_integration_ecr.example EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).


