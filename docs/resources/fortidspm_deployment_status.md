---
subcategory: "Cloud Account Integrations"
layout: "lacework"
page_title: "Lacework: lacework_fortidspm_deployment_status"
description: |-
  Report the outcome of a FortiDSPM scan engine deployment
---

# lacework\_fortidspm\_deployment\_status

Use this resource to tell FortiDSPM how a scan engine deployment ended. The
scan engine modules ([terraform-aws-fortidspm](https://github.com/lacework/terraform-aws-fortidspm),
[terraform-azure-fortidspm](https://github.com/lacework/terraform-azure-fortidspm))
create it for their own region; you only need it directly when you build the
scan engines yourself.

Because its arguments reference the scan engine resources, terraform creates
it after they exist. An apply that fails earlier sends no report. FortiDSPM
merges regions by name, so one resource per region is fine. Destroying the
resource sends nothing; deleting the integration itself tells FortiDSPM the
deployment is gone.

## Example Usage

```hcl
resource "lacework_fortidspm_deployment_status" "us_west_2" {
  intg_guid     = lacework_integration_aws_fortidspm.main.intg_guid
  deployment_id = lacework_integration_aws_fortidspm.main.deployment_id

  region {
    name                  = "us-west-2"
    instance_id           = aws_instance.scan_engine.id
    nat_gateway_public_ip = aws_eip.nat.public_ip
    iam_role_arn          = aws_iam_role.scan_engine.arn
  }
}
```

## Argument Reference

* `intg_guid` - (Required) GUID of the FortiDSPM-managed DSPM integration. Changing this forces a new resource.
* `deployment_id` - (Required) The FortiDSPM deployment id.
* `status` - (Optional) `succeeded` (default), `failed`, `rolled_back`, `deleted_with_resources` or `deleted_orphaned`.
* `region` - (Required) One or more blocks:
  * `name` - (Required) The region.
  * `status` - (Optional) `succeeded` (default) or `failed`.
  * `instance_id` - (Optional) EC2 instance id or Azure VM resource id.
  * `private_ip`, `nat_gateway_public_ip`, `iam_role_arn`, `identity_principal_id` - (Optional) Resource details shown by FortiDSPM.
* `error_phase`, `error_message` - (Optional) Failure details when `status` is not `succeeded`.
