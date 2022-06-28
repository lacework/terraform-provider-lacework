package lacework

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationGcpGkeAuditLog() *schema.Resource {
	return &schema.Resource{
		Create:   resourceLaceworkIntegrationGcpGkeAuditLogCreate,
		Read:     resourceLaceworkIntegrationGcpGkeAuditLogRead,
		Update:   resourceLaceworkIntegrationGcpGkeAuditLogUpdate,
		Delete:   resourceLaceworkIntegrationGcpGkeAuditLogDelete,
		Schema:   gcpGkeAuditLogIntegrationSchema,
		Importer: &schema.ResourceImporter{State: importLaceworkIntegration},
	}
}

var gcpGkeAuditLogIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The integration name.",
	},
	"intg_guid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "The state of the external integration.",
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
				"client_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"private_key_id": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						// @afiune we can't compare this element since our API, for security reasons,
						// does NOT return the private key configured in the Lacework server. So if
						// any other element changed from the credentials then we trigger a diff
						return !d.HasChanges(
							"name", "integration_type", "project_id", "organization_id",
							"subscription", "enabled",
							"credentials.0.client_id",
							"credentials.0.client_email",
						)
					},
				},
				"client_email": {
					Type:     schema.TypeString,
					Required: true,
				},
				"private_key": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						// @afiune we can't compare this element since our API, for security reasons,
						// does NOT return the private key configured in the Lacework server. So if
						// any other element changed from the credentials then we trigger a diff
						return !d.HasChanges(
							"name", "integration_type", "project_id", "organization_id",
							"subscription", "enabled",
							"credentials.0.client_id",
							"credentials.0.client_email",
						)
					},
				},
			},
		},
	},
	"integration_type": {
		Type:     schema.TypeString,
		Required: true,
		StateFunc: func(val interface{}) string {
			return strings.ToUpper(val.(string))
		},
		ValidateFunc: func(value interface{}, key string) ([]string, []error) {
			switch strings.ToUpper(value.(string)) {
			case "PROJECT", "ORGANIZATION":
				return nil, nil
			default:
				return nil, []error{
					fmt.Errorf("%s: can only be either 'PROJECT' or 'ORGANIZATION'", key),
				}
			}
		},
	},
	"organization_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The GCP Organization ID (required when integration_type is organization).",
	},
	"project_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The GCP Project ID.",
	},
	"subscription": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The PubSub Subscription.",
	},
	"is_org": {
		Type:     schema.TypeBool,
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
}

func resourceLaceworkIntegrationGcpGkeAuditLogCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		retries            = d.Get("retries").(int)
		gcpGkeAuditLogData = api.GcpGkeAuditData{
			Credentials: api.GcpGkeAuditCredentials{
				ClientId:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyId: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
			},
			IntegrationType:  strings.ToUpper(d.Get("integration_type").(string)),
			OrganizationId:   d.Get("organization_id").(string),
			ProjectId:        d.Get("project_id").(string),
			SubscriptionName: d.Get("subscription").(string),
		}
	)

	gcpGkeAuditLog := api.NewCloudAccount(d.Get("name").(string),
		api.GcpGkeAuditCloudAccount,
		gcpGkeAuditLogData,
	)

	if !d.Get("enabled").(bool) {
		gcpGkeAuditLog.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s cloud account integration\n", api.GcpGkeAuditCloudAccount.String())
		response, err := lacework.V2.CloudAccounts.Create(gcpGkeAuditLog)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("error creating %s cloud account integration: %s",
						api.GcpGkeAuditCloudAccount.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s cloud account integration. (retrying %d more time(s))\n%s\n",
				api.GcpGkeAuditCloudAccount.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"unable to create %s cloud account integration (retrying %d more time(s))",
				api.GcpGkeAuditCloudAccount.String(), retries,
			))
		}

		cloudAccount := response.Data
		d.SetId(cloudAccount.IntgGuid)
		d.Set("name", cloudAccount.Name)
		d.Set("intg_guid", cloudAccount.IntgGuid)
		d.Set("enabled", cloudAccount.Enabled == 1)

		d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
		d.Set("type_name", cloudAccount.Type) // @afiune should we deprecate?
		d.Set("is_org", cloudAccount.IsOrg == 1)

		log.Printf("[INFO] Created %s cloud account integration with guid: %v\n",
			api.GcpGkeAuditCloudAccount.String(), cloudAccount.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationGcpGkeAuditLogRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s cloud account integration with guid: %v\n", api.GcpGkeAuditCloudAccount.String(), d.Id())
	response, err := lacework.V2.CloudAccounts.GetGcpGkeAudit(d.Id())
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
		creds["client_id"] = response.Data.Data.Credentials.ClientId
		creds["client_email"] = response.Data.Data.Credentials.ClientEmail

		d.Set("credentials", []map[string]string{creds})
		d.Set("integration_type", cloudAccount.Data.IntegrationType)
		d.Set("organization_id", cloudAccount.Data.OrganizationId)
		d.Set("project_id", cloudAccount.Data.ProjectId)
		d.Set("subscription", cloudAccount.Data.SubscriptionName)

		log.Printf("[INFO] Read %s cloud account integration with guid: %v\n",
			api.GcpGkeAuditCloudAccount.String(), cloudAccount.IntgGuid,
		)
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationGcpGkeAuditLogUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework           = meta.(*api.Client)
		gcpGkeAuditLogData = api.GcpGkeAuditData{
			Credentials: api.GcpGkeAuditCredentials{
				ClientId:     d.Get("credentials.0.client_id").(string),
				ClientEmail:  d.Get("credentials.0.client_email").(string),
				PrivateKeyId: d.Get("credentials.0.private_key_id").(string),
				PrivateKey:   d.Get("credentials.0.private_key").(string),
			},
			IntegrationType:  strings.ToUpper(d.Get("integration_type").(string)),
			OrganizationId:   d.Get("organization_id").(string),
			ProjectId:        d.Get("project_id").(string),
			SubscriptionName: d.Get("subscription").(string),
		}
	)

	gcpGkeAuditLog := api.NewCloudAccount(d.Get("name").(string),
		api.GcpGkeAuditCloudAccount,
		gcpGkeAuditLogData,
	)

	if !d.Get("enabled").(bool) {
		gcpGkeAuditLog.Enabled = 0
	}

	gcpGkeAuditLog.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.GcpGkeAuditCloudAccount.String(), gcpGkeAuditLog.IntgGuid)
	response, err := lacework.V2.CloudAccounts.UpdateGcpGkeAudit(gcpGkeAuditLog)
	if err != nil {
		return err
	}

	cloudAccount := response.Data
	d.Set("name", cloudAccount.Name)
	d.Set("intg_guid", cloudAccount.IntgGuid)
	d.Set("enabled", cloudAccount.Enabled == 1)
	d.Set("created_or_updated_time", cloudAccount.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", cloudAccount.CreatedOrUpdatedBy)
	d.Set("type_name", cloudAccount.Type)
	d.Set("is_org", cloudAccount.IsOrg == 1)

	log.Printf("[INFO] Updated %s cloud account integration with guid: %v\n", api.GcpGkeAuditCloudAccount.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationGcpGkeAuditLogDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s cloud account integration with guid: %v\n", api.GcpGkeAuditCloudAccount.String(), d.Id())
	err := lacework.V2.CloudAccounts.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s cloud account integration with guid: %v\n", api.GcpGkeAuditCloudAccount.String(), d.Id())
	return nil
}
