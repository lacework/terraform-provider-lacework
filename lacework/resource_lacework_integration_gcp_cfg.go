package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationGCPCFG() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationGCPCFGCreate,
		Read:   resourceLaceworkIntegrationGCPCFGRead,
		Update: resourceLaceworkIntegrationGCPCFGUpdate,
		Delete: resourceLaceworkIntegrationGCPCFGDelete,
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
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PROJECT",
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch strings.ToUpper(value.(string)) {
					case "PROJECT", "ORGANIZATION":
						return nil, nil
					default:
						return nil, []error{fmt.Errorf("%s: can only be either 'PROJECT' or 'ORGANIZATION'", key)}
					}
				},
			},
			"resource_id": {
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

func resourceLaceworkIntegrationGCPCFGCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		resourceLevel = api.GcpProject
	)

	if d.Get("resource_level").(string) == "ORGANIZATION" {
		resourceLevel = api.GcpOrganization
	}

	data := api.NewGCPIntegrationData(d.Get("name").(string), resourceLevel)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.Data.ID = d.Get("resource_id").(string)
	data.Data.Credentials.ClientId = d.Get("credentials.0.client_id").(string)
	data.Data.Credentials.ClientEmail = d.Get("credentials.0.client_email").(string)
	data.Data.Credentials.PrivateKeyId = d.Get("credentials.0.private_key_id").(string)
	data.Data.Credentials.PrivateKey = d.Get("credentials.0.private_key").(string)

	log.Printf("[INFO] Creating GCP CFG integration with data:\n%+v\n", data)
	response, err := lacework.CreateGCPConfigIntegration(data)
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)
		d.Set("resource_type", integration.Data.IdType)
		d.Set("resource_id", integration.Data.ID)

		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.TypeName)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created GCP CFG integration with guid: %v\n", integration.IntgGuid)
		return nil
	}

	return nil
}

func resourceLaceworkIntegrationGCPCFGRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading GCP CFG integration with guid: %v\n", d.Id())
	response, err := lacework.GetGCPConfigIntegration(d.Id())

	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("resource_type", integration.Data.IdType)
			d.Set("resource_id", integration.Data.ID)

			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)

			log.Printf("[INFO] Read GCP CFG integration with guid: %v\n", integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGCPCFGUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		resourceLevel = api.GcpProject
	)

	if d.Get("resource_level").(string) == "ORGANIZATION" {
		resourceLevel = api.GcpOrganization
	}

	data := api.NewGCPIntegrationData(d.Get("name").(string), resourceLevel)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()
	data.Data.ID = d.Get("resource_id").(string)
	data.Data.Credentials.ClientId = d.Get("credentials.0.client_id").(string)
	data.Data.Credentials.ClientEmail = d.Get("credentials.0.client_email").(string)
	data.Data.Credentials.PrivateKeyId = d.Get("credentials.0.private_key_id").(string)
	data.Data.Credentials.PrivateKey = d.Get("credentials.0.private_key").(string)

	log.Printf("[INFO] Updating GCP CFG integration with data:\n%+v\n", data)
	response, err := lacework.UpdateGCPConfigIntegration(data)
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("resource_type", integration.Data.IdType)
			d.Set("resource_id", integration.Data.ID)

			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)

			log.Printf("[INFO] Updated GCP CFG integration with guid: %v\n", d.Id())
			return nil
		}
	}

	return nil
}

func resourceLaceworkIntegrationGCPCFGDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting GCP CFG integration with guid: %v\n", d.Id())

	_, err := lacework.DeleteGCPConfigIntegration(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted GCP CFG integration with guid: %v\n", d.Id())

	return nil
}
