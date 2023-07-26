package lacework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

func resourceLaceworkAlertProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertProfileCreate,
		Read:   resourceLaceworkAlertProfileRead,
		Update: resourceLaceworkAlertProfileUpdate,
		Delete: resourceLaceworkAlertProfileDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkAlertProfile,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The name of the alert profile",
				ValidateDiagFunc: StringDoesNotHavePrefix("LW_"),
			},
			"extends": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of existing alert profile from which this profile extends",
			},
			"alert": {
				Type:        schema.TypeSet,
				MinItems:    1,
				Required:    true,
				Description: "The list of alert templates",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name that policies can use to refer to this template when generating alerts",
							Required:    true,
						},
						"event_name": {
							Type:        schema.TypeString,
							Description: "The name of the resulting alert",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The summary of the resulting alert",
							Required:    true,
						},
						"subject": {
							Type:        schema.TypeString,
							Description: "A high-level observation of the resulting alert",
							Required:    true,
						},
					},
				},
			},
			"fields": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A field is a declaration of a field to be mapped in from an LQL query",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceLaceworkAlertProfileCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alerts   []api.AlertTemplate
	)
	err := castSchemaSetToArrayOfAlertTemplate(d, "alert", &alerts)
	if err != nil {
		return err
	}

	alertProfile := api.NewAlertProfile(d.Get("name").(string),
		d.Get("extends").(string),
		alerts)

	log.Printf("[INFO] Creating alert profile with data:\n%+v\n", alertProfile)
	response, err := lacework.V2.Alert.Profiles.Create(alertProfile)
	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("fields", setAlertProfileFields(response.Data.Fields))
	log.Printf("[INFO] Created alert profile with id: %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertProfileRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		response api.AlertProfileResponse
	)

	log.Printf("[INFO] Reading alert profile with id: %s\n", d.Id())
	err := lacework.V2.Alert.Profiles.Get(d.Id(), &response)
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.Guid)
	d.Set("id", response.Data.Guid)
	d.Set("extends", response.Data.Extends)
	d.Set("alerts", response.Data.Alerts)
	d.Set("fields", setAlertProfileFields(response.Data.Fields))

	log.Printf("[INFO] Read alert profile with id: %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		alerts   []api.AlertTemplate
	)
	profileID := d.Get("name").(string)

	err := castSchemaSetToArrayOfAlertTemplate(d, "alert", &alerts)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Updating alert profile %s with data:\n%+v\n", profileID, alerts)

	response, err := lacework.V2.Alert.Profiles.Update(profileID, alerts)
	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("fields", setAlertProfileFields(response.Data.Fields))
	log.Printf("[INFO] Updated alert profile with id: %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkAlertProfileDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting alert profile with id: %s\n", d.Id())
	err := lacework.V2.Alert.Profiles.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted alert profile with id: %s\n", d.Id())
	return nil
}

func importLaceworkAlertProfile(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.AlertProfileResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Alert Profile with id: %s\n", d.Id())

	if err := lacework.V2.Alert.Profiles.Get(d.Id(), &response); err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Alert Profile with id '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Alert Profile found with id: %s\n", response.Data.Guid)
	return []*schema.ResourceData{d}, nil
}

func castSchemaSetToArrayOfAlertTemplate(d *schema.ResourceData, attr string, templateList *[]api.AlertTemplate) (err error) {
	var (
		t    api.AlertTemplate
		list []any
	)

	if d.Get(attr) == nil {
		return errors.Errorf("attribute %s not found", attr)
	}

	list = d.Get(attr).(*schema.Set).List()
	for _, item := range list {
		iMap, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("[WARN] unable to cast alert template %v", item)
			continue
		}
		val := sanitizeAlertTemplateKeys(iMap)
		v, err := json.Marshal(val)
		if err != nil {
			return errors.New("failed to marshall alert template attribute")
		}
		err = json.Unmarshal(v, &t)
		if err != nil {
			return errors.New("failed to unmarshall alert template attribute")
		}
		*templateList = append(*templateList, t)
	}
	return
}

func sanitizeAlertTemplateKeys(itemMap map[string]interface{}) map[string]interface{} {
	var newMap = make(map[string]interface{})
	for k, v := range itemMap {
		newKey := strings.Replace(k, "_", "", -1)
		newMap[newKey] = v
	}
	return newMap
}

func setAlertProfileFields(alertFields []api.AlertProfileField) []string {
	var fields []string
	for _, f := range alertFields {
		fields = append(fields, f.Name)
	}
	return fields
}
