package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationOciCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationOciCfgCreate,
		Read:   resourceLaceworkIntegrationOciCfgRead,
		Update: resourceLaceworkIntegrationOciCfgUpdate,
		Delete: resourceLaceworkIntegrationOciCfgDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkCloudAccount,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The number of attempts to create the external integration.",
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fingerprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// we can't compare this element since our API, for security reasons,
								// does NOT return the private key configured in the Lacework server. So if
								// any other element changed from the credentials then we trigger a diff
								return !d.HasChanges(
									"name", "home_region", "tenant_id", "tenant_name",
									"enabled", "user_ocid", "credentials.0.fingerprint",
								)
							},
						},
					},
				},
			},
			"home_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_ocid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_or_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_level": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkIntegrationOciCfgCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		oci      = api.NewCloudAccount(d.Get("name").(string),
			api.OciCfgCloudAccount,
			api.OciCfgData{
				Credentials: api.OciCfgCredentials{
					Fingerprint: d.Get("credentials.0.fingerprint").(string),
					PrivateKey:  d.Get("credentials.0.private_key").(string),
				},
				HomeRegion: d.Get("home_region").(string),
				TenantID:   d.Get("tenant_id").(string),
				TenantName: d.Get("tenant_name").(string),
				UserOCID:   d.Get("user_ocid").(string),
			})
	)

	if !d.Get("enabled").(bool) {
		oci.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.OciCfgCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(oci)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.OciCfgCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.OciCfgCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.OciCfgCloudAccount.String(), retries,
			))
		}

		integration := response.Data
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.Type)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.OciCfgCloudAccount.String(), integration.IntgGuid)
		return nil
	})

}

func resourceLaceworkIntegrationOciCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.OciCfgCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetOciCfg(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	cloudAccount := response.Data
	if cloudAccount.IntgGuid == d.Id() {
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)
		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type)
		d.Set("org_level", cloudAccount.IsOrg == 1)

		creds := make(map[string]string)
		credentials := cloudAccount.Data.Credentials
		creds["fingerprint"] = credentials.Fingerprint
		creds["private_key"] = credentials.PrivateKey
		d.Set("credentials", []map[string]string{creds})
		d.Set("home_region", cloudAccount.Data.HomeRegion)
		d.Set("tenant_id", cloudAccount.Data.TenantID)
		d.Set("tenant_name", cloudAccount.Data.TenantName)
		d.Set("user_ocid", cloudAccount.Data.UserOCID)

		log.Printf("[INFO] Read %s integration with guid: %v\n",
			api.OciCfgCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationOciCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		oci      = api.NewCloudAccount(d.Get("name").(string),
			api.OciCfgCloudAccount,
			api.OciCfgData{
				Credentials: api.OciCfgCredentials{
					Fingerprint: d.Get("credentials.0.fingerprint").(string),
					PrivateKey:  d.Get("credentials.0.private_key").(string),
				},
				HomeRegion: d.Get("home_region").(string),
				TenantID:   d.Get("tenant_id").(string),
				TenantName: d.Get("tenant_name").(string),
				UserOCID:   d.Get("user_ocid").(string),
			})
	)

	if !d.Get("enabled").(bool) {
		oci.Enabled = 0
	}

	oci.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.OciCfgCloudAccount.String(), oci)
	response, err := lacework.V2.CloudAccounts.UpdateOciCfg(oci)
	if err != nil {
		return err
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.OciCfgCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationOciCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.OciCfgCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.OciCfgCloudAccount.String(), d.Id())
	return nil
}
