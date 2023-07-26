package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkDataExportRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkDataExportRuleCreate,
		Read:   resourceLaceworkDataExportRuleRead,
		Update: resourceLaceworkDataExportRuleUpdate,
		Delete: resourceLaceworkDataExportRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkDataExportRule,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the data export rule.",
				Required:    true,
			},
			"integration_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list alert channel ids for the rule to use.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The summary of the data export rule.",
				Optional:    true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
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
		},
	}
}

func resourceLaceworkDataExportRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework   = meta.(*api.Client)
		exportRule = api.DataExportRule{
			Filter: api.DataExportRuleFilter{
				Name:        d.Get("name").(string),
				Description: d.Get("description").(string),
				Enabled:     1,
			},
			Type: "Dataexport",
			IDs:  castAttributeToStringSlice(d, "integration_ids"),
		}
	)

	if !d.Get("enabled").(bool) {
		exportRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Creating data export rule with data:\n%+v\n", exportRule)
	response, err := lacework.V2.DataExportRules.Create(exportRule)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	d.Set("name", response.Data.Filter.Name)
	d.Set("guid", response.Data.ID)
	d.Set("type", response.Data.Type)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.UpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedBy)

	log.Printf("[INFO] Created data export rule with guid %s\n", response.Data.ID)
	return nil
}

func resourceLaceworkDataExportRuleRead(d *schema.ResourceData, meta interface{}) error {
	var lacework = meta.(*api.Client)

	log.Printf("[INFO] Reading data export rule with guid %s\n", d.Id())
	response, err := lacework.V2.DataExportRules.Get(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.ID)
	d.Set("name", response.Data.Filter.Name)
	d.Set("description", response.Data.Filter.Description)
	d.Set("guid", response.Data.ID)
	d.Set("integration_ids", response.Data.IDs)
	d.Set("type", response.Data.Type)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.UpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedBy)

	log.Printf("[INFO] Read data export rule with guid %s\n", response.Data.ID)
	return nil
}

func resourceLaceworkDataExportRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework   = meta.(*api.Client)
		exportRule = api.DataExportRule{
			Filter: api.DataExportRuleFilter{
				Name:        d.Get("name").(string),
				Description: d.Get("description").(string),
				Enabled:     1,
			},
			Type: "Dataexport",
			IDs:  castAttributeToStringSlice(d, "integration_ids"),
		}
	)

	exportRule.ID = d.Id()

	if !d.Get("enabled").(bool) {
		exportRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Updating data export rule with data:\n%+v\n", exportRule)
	response, err := lacework.V2.DataExportRules.Update(exportRule)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	d.Set("name", response.Data.Filter.Name)
	d.Set("type", response.Data.Type)
	d.Set("guid", response.Data.ID)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.UpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedBy)

	log.Printf("[INFO] Updated data export rule with guid %s\n", response.Data.ID)
	return nil
}

func resourceLaceworkDataExportRuleDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting data export rule with guid: %v\n", d.Id())
	err := lacework.V2.DataExportRules.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted data export rule with guid: %v\n", d.Id())
	return nil
}

func importLaceworkDataExportRule(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Data Export Rule with guid: %s\n", d.Id())

	response, err := lacework.V2.DataExportRules.Get(d.Id())
	if err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Data Export Rule with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Data Export Rule found with guid: %s\n", response.Data.ID)
	return []*schema.ResourceData{d}, nil
}
