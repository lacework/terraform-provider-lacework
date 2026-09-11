package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lacework/go-sdk/v2/api"
)

var fortiDspmDeploymentStatuses = []string{
	"succeeded", "failed", "rolled_back", "deleted_with_resources", "deleted_orphaned",
}

// lacework_fortidspm_deployment_status tells FortiDSPM how the scan engine
// deployment ended. Its arguments reference the scan engine resources, so
// terraform creates it only after they exist; an apply that fails earlier
// never reports, which is the agreed behaviour. FortiDSPM merges regions by
// name, so each scan engine module instance reports its own region.
func resourceLaceworkFortiDspmDeploymentStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkFortiDspmDeploymentStatusCreate,
		Read:   resourceLaceworkFortiDspmDeploymentStatusRead,
		Update: resourceLaceworkFortiDspmDeploymentStatusCreate,
		Delete: resourceLaceworkFortiDspmDeploymentStatusDelete,

		Schema: map[string]*schema.Schema{
			"intg_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GUID of the FortiDSPM-managed DSPM integration.",
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FortiDSPM deployment id the integration returned.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "succeeded",
				ValidateFunc: validation.StringInSlice(fortiDspmDeploymentStatuses, false),
				Description:  "Outcome to report: succeeded, failed, rolled_back, deleted_with_resources or deleted_orphaned.",
			},
			"region": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Outcome and resource ids per region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "succeeded",
							ValidateFunc: validation.StringInSlice([]string{"succeeded", "failed"}, false),
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "EC2 instance id or Azure VM resource id of the scan engine.",
						},
						"private_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nat_gateway_public_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"iam_role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"identity_principal_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Principal id of the Azure managed identity attached to the scan engine.",
						},
					},
				},
			},
			"error_phase": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phase that failed (apply, destroy) when status is not succeeded.",
			},
			"error_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLaceworkFortiDspmDeploymentStatusCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	intgGuid := d.Get("intg_guid").(string)

	report := api.FortiDspmDeploymentStatus{
		DeploymentID: d.Get("deployment_id").(string),
		Status:       d.Get("status").(string),
	}
	regionNames := []string{}
	for _, raw := range d.Get("region").([]interface{}) {
		r := raw.(map[string]interface{})
		regionNames = append(regionNames, r["name"].(string))
		report.Regions = append(report.Regions, api.FortiDspmRegionStatus{
			Region:              r["name"].(string),
			Status:              r["status"].(string),
			InstanceID:          r["instance_id"].(string),
			PrivateIP:           r["private_ip"].(string),
			NatGatewayPublicIP:  r["nat_gateway_public_ip"].(string),
			IamRoleArn:          r["iam_role_arn"].(string),
			IdentityPrincipalID: r["identity_principal_id"].(string),
		})
	}
	if msg := d.Get("error_message").(string); msg != "" {
		report.Error = &api.FortiDspmDeploymentError{Phase: d.Get("error_phase").(string), Message: msg}
	}

	log.Printf("[INFO] Reporting FortiDSPM deployment status %s for %s (%s)\n",
		report.Status, intgGuid, strings.Join(regionNames, ","))
	if err := lacework.V2.CloudAccounts.ReportFortiDspmDeploymentStatus(intgGuid, report); err != nil {
		return fmt.Errorf("Error reporting FortiDSPM deployment status for %s: %s", intgGuid, err)
	}
	d.SetId(intgGuid + ":" + strings.Join(regionNames, ","))
	return nil
}

// The report is write-only: FortiDSPM keeps the merged outcome, nothing is read back.
func resourceLaceworkFortiDspmDeploymentStatusRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// Removing the report does not undo it. The integration's own delete tells
// FortiDSPM the deployment is gone.
func resourceLaceworkFortiDspmDeploymentStatusDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
