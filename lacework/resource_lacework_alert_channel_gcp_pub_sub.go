package lacework

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLaceworkAlertChannelGcpPubSub() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelGcpPubSubCreate,
		Read:   resourceLaceworkAlertChannelGcpPubSubRead,
		Update: resourceLaceworkAlertChannelGcpPubSubUpdate,
		Delete: resourceLaceworkAlertChannelGcpPubSubDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
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
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"client_email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceLaceworkAlertChannelGcpPubSubCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		s3       = api.NewGcpPubSubAlertChannel(d.Get("name").(string),
			api.GcpPubSubChannelData{
				ProjectID: d.Get("project_id").(string),
				TopicID:   d.Get("topic_id").(string),
				Credentials: api.GcpCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientEmail:  d.Get("credentials.0.client_email").(string),
					PrivateKey:   d.Get("credentials.0.private_key").(string),
					PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		s3.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.GcpPubSubChannelIntegration, s3)
	response, err := lacework.Integrations.CreateGcpPubSubAlertChannel(s3)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateGcpPubSubAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	integration := response.Data[0]
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Created %s integration with guid: %v\n", api.GcpPubSubChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n", api.GcpPubSubChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetGcpPubSubAlertChannel(d.Id())
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)
			d.Set("project_id", integration.Data.ProjectID)
			d.Set("topic_id", integration.Data.TopicID)

			creds := make(map[string]string)
			creds["client_id"] = integration.Data.Credentials.ClientID
			creds["client_email"] = integration.Data.Credentials.ClientEmail
			creds["private_key"] = integration.Data.Credentials.PrivateKey
			creds["private_key_id"] = integration.Data.Credentials.PrivateKeyID

			d.Set("credentials", []map[string]string{creds})

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.GcpPubSubChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		s3       = api.NewGcpPubSubAlertChannel(d.Get("name").(string),
			api.GcpPubSubChannelData{
				ProjectID: d.Get("project_id").(string),
				TopicID:   d.Get("topic_id").(string),
				Credentials: api.GcpCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientEmail:  d.Get("credentials.0.client_email").(string),
					PrivateKey:   d.Get("credentials.0.private_key").(string),
					PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		s3.Enabled = 0
	}

	s3.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.GcpPubSubChannelIntegration, s3)
	response, err := lacework.Integrations.UpdateGcpPubSubAlertChannel(s3)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateGcpPubSubAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n", api.GcpPubSubChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n", api.GcpPubSubChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n", api.GcpPubSubChannelIntegration, d.Id())
	return nil
}

func validateGcpPubSubAlertChannelResponse(response *api.GcpPubSubAlertChannelResponse) error {
	if len(response.Data) == 0 {
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		msg := `
There is more that one integration inside the server response data.

List of integrations:
`
		for _, integration := range response.Data {
			msg = msg + fmt.Sprintf("\t%s: %s\n", integration.IntgGuid, integration.Name)
		}
		msg = msg + unexpectedBehaviorMsg()
		return fmt.Errorf(msg)
	}

	return nil
}
