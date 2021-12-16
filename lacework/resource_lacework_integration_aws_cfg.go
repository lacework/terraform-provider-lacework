package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkIntegrationAwsCfg() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationAwsCfgCreate,
		Read:   resourceLaceworkIntegrationAwsCfgRead,
		Update: resourceLaceworkIntegrationAwsCfgUpdate,
		Delete: resourceLaceworkIntegrationAwsCfgDelete,

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
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceLaceworkIntegrationAwsCfgCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		retries  = d.Get("retries").(int)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCfgIntegration,
			api.AwsIntegrationData{
				Credentials: &api.AwsCrossAccountCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		retries--
		log.Printf("[INFO] Creating %s integration\n", api.AwsCfgIntegration.String())
		response, err := lacework.Integrations.CreateAws(aws)
		if err != nil {
			if retries <= 0 {
				return resource.NonRetryableError(
					fmt.Errorf("Error creating %s integration: %s",
						api.AwsCfgIntegration.String(), err,
					))
			}
			log.Printf(
				"[INFO] Unable to create %s integration. (retrying %d more time(s))\n%s\n",
				api.AwsCfgIntegration.String(), retries, err,
			)
			return resource.RetryableError(fmt.Errorf(
				"Unable to create %s integration (retrying %d more time(s))",
				api.AwsCfgIntegration.String(), retries,
			))
		}

		log.Printf("[INFO] Verifying server response")
		err = validateAwsIntegrationResponse(&response)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		// @afiune at this point in time, we know the data field has a single value
		integration := response.Data[0]
		d.SetId(integration.IntgGuid)
		d.Set("name", integration.Name)
		d.Set("intg_guid", integration.IntgGuid)
		d.Set("enabled", integration.Enabled == 1)

		d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
		d.Set("type_name", integration.TypeName)
		d.Set("org_level", integration.IsOrg == 1)

		log.Printf("[INFO] Created %s integration with guid: %v\n",
			api.AwsCfgIntegration.String(), integration.IntgGuid)
		return nil
	})
}

func resourceLaceworkIntegrationAwsCfgRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid: %v\n",
		api.AwsCfgIntegration.String(), d.Id())
	response, err := lacework.Integrations.GetAws(d.Id())
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

			creds := make(map[string]string)
			credentials := integration.Data.GetCredentials()
			creds["role_arn"] = credentials.RoleArn
			creds["external_id"] = credentials.ExternalID
			d.Set("credentials", []map[string]string{creds})

			log.Printf("[INFO] Read %s integration with guid: %v\n",
				api.AwsCfgIntegration.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationAwsCfgUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		aws      = api.NewAwsIntegration(d.Get("name").(string),
			api.AwsCfgIntegration,
			api.AwsIntegrationData{
				Credentials: &api.AwsCrossAccountCreds{
					RoleArn:    d.Get("credentials.0.role_arn").(string),
					ExternalID: d.Get("credentials.0.external_id").(string),
				},
			},
		)
	)

	if !d.Get("enabled").(bool) {
		aws.Enabled = 0
	}

	aws.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n",
		api.AwsCfgIntegration.String(), aws)
	response, err := lacework.Integrations.UpdateAws(aws)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAwsIntegrationResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point in time, we know the data field has a single value
	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	log.Printf("[INFO] Updated %s integration with guid: %v\n",
		api.AwsCfgIntegration.String(), d.Id())
	return nil
}

func resourceLaceworkIntegrationAwsCfgDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid: %v\n",
		api.AwsCfgIntegration.String(), d.Id())
	_, err := lacework.Integrations.DeleteAws(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid: %v\n",
		api.AwsCfgIntegration.String(), d.Id())
	return nil
}

// validateAwsIntegrationResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateAwsIntegrationResponse(response *api.AwsIntegrationsResponse) error {
	if len(response.Data) == 0 {
		// @afiune this edge case should never happen, if we land here it means that
		// something went wrong in the server side of things (Lacework API), so let
		// us inform that to our users
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		// @afiune if we are creating a single integration and the server returns
		// more than one integration inside the 'data' field, it is definitely another
		// edge case that should never happen
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

func unexpectedBehaviorMsg() string {
	return `
This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
}
