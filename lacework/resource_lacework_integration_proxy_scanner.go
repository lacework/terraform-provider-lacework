package lacework

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationProxyScanner() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationProxyScannerCreate,
		Read:   resourceLaceworkIntegrationProxyScannerRead,
		Update: resourceLaceworkIntegrationProxyScannerUpdate,
		Delete: resourceLaceworkIntegrationProxyScannerDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkContainerRegistry,
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
			"limit_by_tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A list of image tags to limit the assessment of images with matching tags",
			},
			"limit_by_label": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},

						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
						},
					},
				},
				Optional:    true,
				Description: "A list of key/value labels to limit the assessment of images",
			},
			"limit_by_repositories": {
				Type:     schema.TypeList,
				MinItems: 0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A list of repositories to assess",
			},
			"limit_num_imgs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
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
			"policy_evaluate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"policy_guids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceLaceworkIntegrationProxyScannerCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(
		d.Get("name").(string),
		api.ProxyScannerContainerRegistry,
		api.ProxyScannerData{
			LimitByTag:   castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel: castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:   castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:  d.Get("limit_num_imgs").(int),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Creating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.ProxyScannerContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.Create(
		data,
	)
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
	d.Set("server_token", integration.ServerToken.ServerToken)
	d.Set("server_token_uri", integration.ServerToken.Uri)

	log.Printf("[INFO] Created ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), integration.IntgGuid)

	if d.Get("policy_evaluate").(bool) {
		log.Printf("[INFO] Map policies...\n")
		_, err := lacework.V2.ContainerRegistries.MapPolicy(
			response.Data.IntgGuid,
			api.MapPolicyRequest{
				Evaluate:    d.Get("policy_evaluate").(bool),
				PolicyGuids: castAttributeToStringSlice(d, "policy_guids"),
			},
		)
		if err != nil {
			return err
		}
	} else {
		d.Set("policy_guids", nil)
		d.Set("policy_evaluate", false)
	}

	return nil
}

func resourceLaceworkIntegrationProxyScannerRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), d.Id())
	response, err := lacework.V2.ContainerRegistries.GetProxyScanner(d.Id())
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
	d.Set("server_token_uri", integration.ServerToken.Uri)

	// check for props and marshal
	if t, ok := integration.Props.(map[string]interface{}); ok {
		if jsonbody, err := json.Marshal(t); err != nil {
			return err
		} else {
			props := api.V2IntegrationProps{}
			if err := json.Unmarshal(jsonbody, &props); err != nil {
				return err
			}
			nop := props.PolicyEvaluation
			if nop != nil {
				log.Printf("[INFO] Found inline policy evaluation: %s\n", strconv.FormatBool(nop.Evaluate))
				d.Set("policy_evaluate", nop.Evaluate)
				if nop.Evaluate {
					for _, nog := range nop.PolicyGuids {
						log.Printf("[INFO] Found inline policy guid: %s\n", nog)
					}
					d.Set("policy_guids", nop.PolicyGuids)
				}
			}
		}
	} else {
		d.Set("policy_guids", nil)
		d.Set("policy_evaluate", false)
	}

	d.Set("limit_num_imgs", integration.Data.LimitNumImg)
	d.Set("limit_by_tags", integration.Data.LimitByTag)
	d.Set("limit_by_repositories", integration.Data.LimitByRep)
	d.Set("limit_by_label", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(integration.Data.LimitByLabel))

	log.Printf("[INFO] Read ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationProxyScannerUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.ProxyScannerContainerRegistry,
		api.ProxyScannerData{
			LimitByTag:   castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel: castAttributeToArrayOfKeyValueMap(d, "limit_by_label"),
			LimitByRep:   castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:  d.Get("limit_num_imgs").(int),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating ContVulnCfg integration for %s registry type with data:\n%+v\n",
		api.ProxyScannerContainerRegistry.String(), data)
	response, err := lacework.V2.ContainerRegistries.UpdateProxyScanner(data)
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
	d.Set("server_token_uri", integration.ServerToken.Uri)

	if d.Get("policy_evaluate").(bool) {
		log.Printf("[INFO] Map policies...\n")
		_, err := lacework.V2.ContainerRegistries.MapPolicy(
			response.Data.IntgGuid,
			api.MapPolicyRequest{
				Evaluate:    d.Get("policy_evaluate").(bool),
				PolicyGuids: castAttributeToStringSlice(d, "policy_guids"),
			},
		)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[INFO] Unmap policies...\n")
		_, err := lacework.V2.ContainerRegistries.MapPolicy(
			response.Data.IntgGuid,
			api.MapPolicyRequest{
				Evaluate:    d.Get("policy_evaluate").(bool),
				PolicyGuids: []string{},
			},
		)
		if err != nil {
			return err
		}
		d.Set("policy_guids", nil)
		d.Set("policy_evaluate", false)
	}

	d.Set("limit_num_imgs", integration.Data.LimitNumImg)
	d.Set("limit_by_tags", integration.Data.LimitByTag)
	d.Set("limit_by_repositories", integration.Data.LimitByRep)
	d.Set("limit_by_label", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(integration.Data.LimitByLabel))

	log.Printf("[INFO] Updated ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationProxyScannerDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), d.Id())
	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted ContVulnCfg integration for %s registry type with guid %s\n",
		api.ProxyScannerContainerRegistry.String(), d.Id())
	return nil
}
