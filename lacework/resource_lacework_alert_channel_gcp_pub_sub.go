package lacework

import (
	"fmt"
	"log"

	"github.com/lacework/go-sdk/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLaceworkAlertChannelGcpPubSub() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelGcpPubSubCreate,
		Read:   resourceLaceworkAlertChannelGcpPubSubRead,
		Update: resourceLaceworkAlertChannelGcpPubSubUpdate,
		Delete: resourceLaceworkAlertChannelGcpPubSubDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAlertChannel,
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
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Google Cloud Project",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Google Cloud Pub/Sub topic",
			},
			"issue_grouping": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Events",
				Description: "Defines how Lacework compliance events get grouped. Must be one of Events or Resources. Defaults to Events.",
				ValidateFunc: func(value interface{}, key string) ([]string, []error) {
					switch value.(string) {
					case "Events", "Resources":
						return nil, nil
					default:
						return nil, []error{
							fmt.Errorf(
								"%s: can only be either 'Events' or 'Resources' (default: Events)", key,
							),
						}
					}
				},
			},
			"credentials": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The service account client ID",
						},
						"client_email": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The service account client email",
						},
						"private_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The service account private key ID",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return !d.HasChanges(
									"name", "project_id", "topic_id", "org_level", "enabled", "issue_grouping",
									"credentials.0.client_id",
									"credentials.0.client_email",
								)
							},
						},
						"private_key_id": {
							Type:        schema.TypeString,
							Sensitive: 	 true,
							Required:    true,
							Description: "The service account private key",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return !d.HasChanges(
									"name", "project_id", "topic_id", "org_level", "enabled", "issue_grouping",
									"credentials.0.client_id",
									"credentials.0.client_email",
								)
							},
						},
					},
				},
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modification",
			},
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intg_guid": {
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

func resourceLaceworkAlertChannelGcpPubSubCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework  = meta.(*api.Client)
		gcpPubSub = api.NewAlertChannel(d.Get("name").(string),
			api.GcpPubSubAlertChannelType,
			api.GcpPubSubDataV2{
				ProjectID:     d.Get("project_id").(string),
				TopicID:       d.Get("topic_id").(string),
				IssueGrouping: d.Get("issue_grouping").(string),
				Credentials: api.GcpPubSubCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientEmail:  d.Get("credentials.0.client_email").(string),
					PrivateKey:   d.Get("credentials.0.private_key").(string),
					PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				},
			},
		)
	)
	if !d.Get("enabled").(bool) {
		gcpPubSub.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.GcpPubSubAlertChannelType, gcpPubSub)
	response, err := lacework.V2.AlertChannels.Create(gcpPubSub)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
		if err := VerifyAlertChannelAndRollback(d, lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.GcpPubSubAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.GcpPubSubAlertChannelType, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
	response, err := lacework.V2.AlertChannels.GetGcpPubSub(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", response.Data.Name)
	d.Set("intg_guid", response.Data.IntgGuid)
	d.Set("enabled", response.Data.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("org_level", response.Data.IsOrg == 1)
	d.Set("project_id", response.Data.Data.ProjectID)
	d.Set("topic_id", response.Data.Data.TopicID)
	d.Set("issue_grouoing", response.Data.Data.IssueGrouping)

	creds := make(map[string]string)
	creds["client_id"] = response.Data.Data.Credentials.ClientID
	creds["client_email"] = response.Data.Data.Credentials.ClientEmail
	creds["private_key"] = response.Data.Data.Credentials.PrivateKey
	creds["private_key_id"] = response.Data.Data.Credentials.PrivateKeyID

	d.Set("credentials", []map[string]string{creds})

	log.Printf("[INFO] Read %s integration with guid %s\n",
		api.GcpPubSubAlertChannelType, response.Data.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework  = meta.(*api.Client)
		gcpPubSub = api.NewAlertChannel(d.Get("name").(string),
			api.GcpPubSubAlertChannelType,
			api.GcpPubSubDataV2{
				ProjectID:     d.Get("project_id").(string),
				TopicID:       d.Get("topic_id").(string),
				IssueGrouping: d.Get("issue_grouping").(string),
				Credentials: api.GcpPubSubCredentials{
					ClientID:     d.Get("credentials.0.client_id").(string),
					ClientEmail:  d.Get("credentials.0.client_email").(string),
					PrivateKey:   d.Get("credentials.0.private_key").(string),
					PrivateKeyID: d.Get("credentials.0.private_key_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		gcpPubSub.Enabled = 0
	}

	gcpPubSub.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.GcpPubSubAlertChannelType, gcpPubSub)
	response, err := lacework.V2.AlertChannels.UpdateGcpPubSub(gcpPubSub)
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

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.GcpPubSubAlertChannelType, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
	return nil
}

func resourceLaceworkAlertChannelGcpPubSubDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
	err := lacework.V2.AlertChannels.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.GcpPubSubAlertChannelType, d.Id())
	return nil
}
