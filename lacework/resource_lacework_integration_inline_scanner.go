package lacework

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/v2/api"
)

func resourceLaceworkIntegrationInlineScanner() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationInlineScannerCreate,
		Read:   resourceLaceworkIntegrationInlineScannerRead,
		Update: resourceLaceworkIntegrationInlineScannerUpdate,
		Delete: resourceLaceworkIntegrationInlineScannerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The integration name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"identifier_tag": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional:    true,
				Description: "A list of identifier tags as key:value pairs",
			},
			"limit_num_scan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The maximum number of newest container images to assess per repository.",
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
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
			"server_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkIntegrationInlineScannerCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(
		d.Get("name").(string),
		api.InlineScannerContainerRegistry,
		api.InlineScannerData{
			IdentifierTag: castAttributeToArrayOfKeyValueMap(d, "identifier_tag"),
			LimitNumScan:  strconv.Itoa(d.Get("limit_num_scan").(int)),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.InlineScannerContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.Create(data)
	if err != nil {
		return err
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
	d.Set("server_token", integration.ServerToken.ServerToken)
	d.Set("server_uri", integration.ServerToken.Uri)

	log.Printf("[INFO] Created ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), response.Data.IntgGuid)

	return nil
}

func resourceLaceworkIntegrationInlineScannerRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetInlineScanner(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)
	d.Set("server_token", integration.ServerToken.ServerToken)
	d.Set("server_uri", integration.ServerToken.Uri)

	if limit_num_scan, err := strconv.Atoi(integration.Data.LimitNumScan); err == nil {
		log.Printf("[INFO] Setting limit_num_scan: %d\n", limit_num_scan)
		d.Set("limit_num_scan", limit_num_scan)
	}
	d.Set("identifier_tag", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(integration.Data.IdentifierTag))

	log.Printf("[INFO] Read ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationInlineScannerUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.InlineScannerContainerRegistry,
		api.InlineScannerData{
			IdentifierTag: castAttributeToArrayOfKeyValueMap(d, "identifier_tag"),
			LimitNumScan:  strconv.Itoa(d.Get("limit_num_scan").(int)),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.InlineScannerContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateInlineScanner(data)
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
	d.Set("server_token", integration.ServerToken.ServerToken)
	d.Set("server_uri", integration.ServerToken.Uri)

	if limit_num_scan, err := strconv.Atoi(integration.Data.LimitNumScan); err == nil {
		log.Printf("[INFO] Setting limit_num_scan: %d\n", limit_num_scan)
		d.Set("limit_num_scan", limit_num_scan)
	}
	d.Set("identifier_tag", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(integration.Data.IdentifierTag))

	log.Printf("[INFO] Updated ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationInlineScannerDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), d.Id())
	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted ContVulnCfg integration for %s registry type with guid %s\n",
		api.InlineScannerContainerRegistry.String(), d.Id())
	return nil
}
