---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_integration_azure_fortidspm"
description: |-
  Create and manage a FortiDSPM-managed Azure DSPM integration
---

# lacework\_integration\_azure\_fortidspm

Use this resource to register an Azure tenant or subscription for DSPM scanning
by FortiDSPM.

Creating the integration also creates the FortiDSPM deployment for it. FortiDSPM
returns, per region, one single-use activation token and a signed URL to the
scan engine image VHD together with its Hyper-V generation. They are exposed as
`activation_tokens`, `image_urls` and `hyperv_generations`, keyed by region, so
the [terraform-azure-fortidspm](https://github.com/lacework/terraform-azure-fortidspm)
module can copy the image into the subscription and boot the scan engines.

The deployment attributes are issued once on create and are not refreshed on
read. Changing `tenant_id`, `subscription_id` or `regions` recreates the
integration and the deployment.

## Example Usage

The scan engine module creates this resource itself when `global = true`, so
the usual way to use it is one module per region:

```hcl
module "lacework_azure_fortidspm_westus2" {
  source                    = "git::https://github.com/lacework/terraform-azure-fortidspm.git?ref=v0.2.0"
  global                    = true
  lacework_integration_name = "azure-dspm-production"
  tenant_id                 = "00000000-0000-0000-0000-000000000000"
  subscription_id           = "11111111-1111-1111-1111-111111111111"
  regions                   = ["westus2", "eastus"]
  location                  = "westus2"
}

module "lacework_azure_fortidspm_eastus" {
  source                  = "git::https://github.com/lacework/terraform-azure-fortidspm.git?ref=v0.2.0"
  global_module_reference = module.lacework_azure_fortidspm_westus2
  location                = "eastus"
}
```

Used directly, the resource registers the tenant and hands out the
per-region inputs:

```hcl
resource "lacework_integration_azure_fortidspm" "main" {
  name            = "azure-dspm-production"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  subscription_id = "11111111-1111-1111-1111-111111111111"
  regions         = ["westus2"]
}

output "image_url_westus2" {
  value     = lacework_integration_azure_fortidspm.main.image_urls["westus2"]
  sensitive = true
}
```

## Argument Reference

* `name` - (Required) The integration name.
* `tenant_id` - (Required) The Azure tenant where the scan engines are deployed. Changing this forces a new resource.
* `subscription_id` - (Optional) The Azure subscription where the scan engines are deployed. Leave empty for a tenant-level integration. Changing this forces a new resource.
* `regions` - (Required) The Azure locations where a scan engine is deployed. Changing this forces a new resource.
* `retries` - (Optional) The number of attempts to create the integration. Defaults to `5`.

## Attribute Reference

* `intg_guid` - The integration GUID.
* `deployment_id` - The FortiDSPM deployment id.
* `deployment_name` - The FortiDSPM deployment name.
* `env_id` - The FortiDSPM environment id.
* `token_expires_in` - Seconds the activation tokens stay valid after creation.
* `image_url_expires_in` - Seconds the image URLs stay valid after creation.
* `activation_tokens` - Map of region to single-use activation token. Sensitive.
* `image_urls` - Map of region to signed scan engine image VHD URL. Sensitive.
* `hyperv_generations` - Map of region to the image's Hyper-V generation (`V1` or `V2`).

## Import

A FortiDSPM-managed Azure DSPM integration can be imported using its `INT_GUID`:

```
$ terraform import lacework_integration_azure_fortidspm.main EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```

The deployment attributes are not available after an import.
