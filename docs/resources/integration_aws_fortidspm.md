---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_aws_fortidspm"
description: |-
  Create and manage a FortiDSPM-managed AWS DSPM integration
---

# lacework\_integration\_aws\_fortidspm

Use this resource to register an AWS account for DSPM scanning by FortiDSPM.

Creating the integration also creates the FortiDSPM deployment for it. FortiDSPM
returns one single-use activation token and one scan engine AMI per region;
they are exposed as `activation_tokens` and `image_ids`, keyed by region, so the
[terraform-aws-fortidspm](https://github.com/lacework/terraform-aws-fortidspm)
module can boot the scan engines. Unlike `lacework_integration_aws_dspm`, no
result bucket and no cross-account role are needed: FortiDSPM pushes results to
FortiCNAPP itself.

The deployment attributes are issued once on create and are not refreshed on
read. Changing `account_id` or `regions` recreates the integration and the
deployment.

## Example Usage

The scan engine module creates this resource itself when `global = true`, so
the usual way to use it is one module per region:

```hcl
provider "aws" {
  alias  = "us_west_2"
  region = "us-west-2"
}

provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

module "lacework_aws_fortidspm_us_west_2" {
  source                    = "git::https://github.com/lacework/terraform-aws-fortidspm.git?ref=v0.2.0"
  global                    = true
  lacework_integration_name = "aws-dspm-123456789012"
  regions                   = ["us-west-2", "us-east-1"]

  providers = { aws = aws.us_west_2 }
}

module "lacework_aws_fortidspm_us_east_1" {
  source                  = "git::https://github.com/lacework/terraform-aws-fortidspm.git?ref=v0.2.0"
  global_module_reference = module.lacework_aws_fortidspm_us_west_2

  providers = { aws = aws.us_east_1 }
}
```

Used directly, the resource registers the account and hands out the
per-region inputs:

```hcl
resource "lacework_integration_aws_fortidspm" "main" {
  name       = "aws-dspm-123456789012"
  account_id = "123456789012"
  regions    = ["us-west-2", "us-east-1"]
}

output "activation_token_us_west_2" {
  value     = lacework_integration_aws_fortidspm.main.activation_tokens["us-west-2"]
  sensitive = true
}
```

## Argument Reference

* `name` - (Required) The integration name.
* `account_id` - (Required) The AWS account where the scan engines are deployed. Changing this forces a new resource.
* `regions` - (Required) The regions where a scan engine is deployed. Changing this forces a new resource.
* `retries` - (Optional) The number of attempts to create the integration. Defaults to `5`.

## Attribute Reference

* `intg_guid` - The integration GUID.
* `deployment_id` - The FortiDSPM deployment id.
* `deployment_name` - The FortiDSPM deployment name.
* `env_id` - The FortiDSPM environment id.
* `token_expires_in` - Seconds the activation tokens stay valid after creation.
* `activation_tokens` - Map of region to single-use activation token. Sensitive.
* `image_ids` - Map of region to scan engine AMI id.

## Import

A FortiDSPM-managed AWS DSPM integration can be imported using its `INT_GUID`:

```
$ terraform import lacework_integration_aws_fortidspm.main EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```

The deployment attributes are not available after an import.
