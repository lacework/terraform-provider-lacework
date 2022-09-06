package lacework

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

func importLaceworkECRIntegration(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)
	var awsAuthType string

	log.Printf("[INFO] Importing Lacework integration with guid: %s\n", d.Id())

	var response api.ContainerRegistryRaw

	err := lacework.V2.ContainerRegistries.Get(d.Id(), &response)
	if err != nil {
		return nil, err
	}

	integration := response.IntgGuid
	if integration == d.Id() {
		log.Printf("[INFO] ECR integration found with guid: %v\n", integration)

		awsAuthType, err = getAuthTypeFromRaw(response)
		if err != nil || awsAuthType == "" {
			return nil, errors.Wrapf(err, "unable to import Lacework resource with guid '%s'",
				d.Id())
		}

		d.Set("aws_auth_type", awsAuthType)

		return []*schema.ResourceData{d}, nil
	}

	log.Printf("[INFO] Raw integration response: %v\n", response)
	return nil, fmt.Errorf(
		"unable to import Lacework resource. Integration with guid '%s' was not found",
		d.Id(),
	)
}

func getAuthTypeFromRaw(raw api.ContainerRegistryRaw) (string, error) {
	if casting, ok := raw.GetData().(map[string]interface{}); ok {
		if _, exist := casting["data"]; exist {
			if castData, ok := casting["data"].(map[string]interface{}); ok {
				if t, exist := castData["awsAuthType"]; exist {
					return t.(string), nil
				}
			}
		}

	}

	return "", errors.New("field AwsAuthType not found in response.")

}

func resourceLaceworkIntegrationEcr() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkIntegrationEcrCreate,
		Read:   resourceLaceworkIntegrationEcrRead,
		Update: resourceLaceworkIntegrationEcrUpdate,
		Delete: resourceLaceworkIntegrationEcrDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkECRIntegration,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ECR integration name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the external integration",
			},
			"registry_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Amazon Container Registry (ECR) domain",
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "The credentials needed by the integration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ARN of the IAM role with permissions to access the Amazon Container Registry",
						},
						"external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The external ID for the IAM role",
						},
						"access_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The AWS access key ID for an AWS IAM user that permissions to access the Amazon Container Registry",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The AWS secret key for the specified AWS access key",
						},
					},
				},
			},
			"non_os_package_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable program language scanning",
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
			"limit_by_labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:    true,
				Description: "A key based map of labels to limit the assessment of images with matching key:value labels",
			},
			"limit_by_repositories": {
				Type: schema.TypeList,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntInSlice([]int{5, 10, 15}),
				Description:  "The maximum number of newest container images to assess per repository",
			},
			"aws_auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authentication method of the ECR integration",
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
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not this integration is configured at the Organization level",
			},
		},
	}
}

func resourceLaceworkIntegrationEcrCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	switch detectAuthenticationMethod(d) {
	case api.AwsEcrAccessKey.String():
		return resourceLaceworkIntegrationEcrCreateWithAccessKey(d, lacework)

	case api.AwsEcrIAM.String():
		return resourceLaceworkIntegrationEcrCreateWithIAMRole(d, lacework)

	default:
		msg := `Invalid credentials block.

For AWS IAM Role-Based Authentication, provide only the arguments:
  * role_arn
  * external_id

For AWS Access Key-Based Authentication, provide only the arguments:
  * access_key_id
  * secret_access_key

For more information visit https://registry.terraform.io/providers/lacework/lacework/latest/docs/resources/integration_ecr
`
		return errors.New(msg)
	}
}

func resourceLaceworkIntegrationEcrRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), d.Id())

	switch d.Get("aws_auth_type").(string) {
	case api.AwsEcrIAM.String():
		return readEcrIam(d, meta)
	case api.AwsEcrAccessKey.String():
		return readEcrAccessKey(d, meta)
	default:
		return errors.Errorf("Unsupported ECR authentication type '%s'.", d.Get("aws_auth_type").(string))
	}
}

func resourceLaceworkIntegrationEcrUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	switch d.Get("aws_auth_type").(string) {
	case api.AwsEcrAccessKey.String():
		log.Printf("[INFO] %s Authentication: %s\n", api.EcrRegistry.String(), api.AwsEcrAccessKey.String())

		if err := validateAccessKeyCreds(d); err != nil {
			// verify if the user is trying to change the authentication method
			errIAM := validateIAMRoleCreds(d)
			if errIAM != nil {
				// nope, just throw the first error as usual
				return err
			}

			// yup, the user is trying to change the authentication method
			// we need to delete the integration and create a new one
			log.Println("[WARN] Change of authentication method detected. Need destroy and recreation")
			log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())
			_, err := lacework.Integrations.Delete(d.Id())
			if err != nil {
				return err
			}

			return resourceLaceworkIntegrationEcrCreate(d, meta)
		}

		return resourceLaceworkIntegrationEcrUpdateWithAccessKey(d, lacework)

	case api.AwsEcrIAM.String():
		log.Printf("[INFO] %s Authentication: %s\n", api.EcrRegistry.String(), api.AwsEcrIAM.String())

		if err := validateIAMRoleCreds(d); err != nil {
			// verify if the user is trying to change the authentication method
			errKeys := validateAccessKeyCreds(d)
			if errKeys != nil {
				// nope, just throw the first error as usual
				return err
			}

			// yup, the user is trying to change the authentication method
			// we need to delete the integration and create a new one
			log.Println("[WARN] Change of authentication method detected. Need destroy and recreation")
			log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())
			_, err := lacework.Integrations.Delete(d.Id())
			if err != nil {
				return err
			}

			return resourceLaceworkIntegrationEcrCreate(d, meta)
		}

		return resourceLaceworkIntegrationEcrUpdateWithIAMRole(d, lacework)

	default:
		return errors.Errorf("Unsupported ECR authentication type '%s'.", d.Get("aws_auth_type").(string))
	}
}

func resourceLaceworkIntegrationEcrUpdateWithIAMRole(d *schema.ResourceData, lacework *api.Client) error {
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.AwsEcrContainerRegistry,
		api.AwsEcrIamRoleData{
			CrossAccountCredentials: api.AwsEcrCrossAccountCredentials{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s registry\n", api.AwsEcrContainerRegistry.String())
	response, err := lacework.V2.ContainerRegistries.UpdateAwsEcrIamRole(data)
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
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Updated  %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), d.Id())

	return nil
}

func resourceLaceworkIntegrationEcrUpdateWithAccessKey(d *schema.ResourceData, lacework *api.Client) error {
	data := api.NewContainerRegistry(d.Get("name").(string),
		api.AwsEcrContainerRegistry,
		api.AwsEcrAccessKeyData{
			AccessKeyCredentials: api.AwsEcrAccessKeyCredentials{
				AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
				SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
			},
			LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
			LimitByLabel:     castAttributeToArrayOfKeyValueMap(d, "limit_by_labels"),
			LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
			LimitNumImg:      d.Get("limit_num_imgs").(int),
			RegistryDomain:   d.Get("registry_domain").(string),
			NonOSPackageEval: d.Get("non_os_package_support").(bool),
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s registry\n", api.AwsEcrContainerRegistry.String())
	response, err := lacework.V2.ContainerRegistries.UpdateAwsEcrAccessKey(data)
	if err != nil {
		return err
	}

	// @afiune at this point in time, we know the data field has a single value
	integration := response.Data
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.Type)
	d.Set("org_level", integration.IsOrg == 1)
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Updated %s registry type with guid: %v\n", api.EcrRegistry.String(), d.Id())

	return nil
}

func resourceLaceworkIntegrationEcrDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), d.Id())

	err := lacework.V2.ContainerRegistries.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), d.Id())

	return nil
}

func detectAuthenticationMethod(d *schema.ResourceData) string {
	if err := validateIAMRoleCreds(d); err == nil {
		return api.AwsEcrIAM.String()
	}

	if err := validateAccessKeyCreds(d); err == nil {
		return api.AwsEcrAccessKey.String()
	}

	return "credentials block misconfigured"
}

func validateIAMRoleCreds(d *schema.ResourceData) error {
	if d.Get("credentials.0.role_arn").(string) == "" {
		return errors.New("missing role_arn argument inside credentials block.")
	}
	if d.Get("credentials.0.external_id").(string) == "" {
		return errors.New("missing external_id argument inside credentials block.")
	}

	return nil
}

func validateAccessKeyCreds(d *schema.ResourceData) error {
	if d.Get("credentials.0.access_key_id").(string) == "" {
		return errors.New("missing access_key_id argument inside credentials block.")
	}
	if d.Get("credentials.0.secret_access_key").(string) == "" {
		return errors.New("missing secret_access_key argument inside credentials block.")
	}

	return nil
}

func resourceLaceworkIntegrationEcrCreateWithAccessKey(d *schema.ResourceData, lacework *api.Client) error {
	iamRoleData := api.AwsEcrAccessKeyData{
		AccessKeyCredentials: api.AwsEcrAccessKeyCredentials{
			AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
			SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
		},
		LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
		LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
		LimitNumImg:      d.Get("limit_num_imgs").(int),
		RegistryDomain:   d.Get("registry_domain").(string),
		NonOSPackageEval: d.Get("non_os_package_support").(bool),
	}

	labels := castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels")
	if len(labels) != 0 {
		iamRoleData.LimitByLabel = labels
	}

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.AwsEcrContainerRegistry,
		iamRoleData,
	)

	log.Printf("[INFO] Creating %s registry\n", api.AwsEcrContainerRegistry.String())
	response, err := lacework.V2.ContainerRegistries.Create(data)
	if err != nil {
		return err
	}

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
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

	ecrData, err := castRawToAwsEcrAccessKeyData(response.Data.Data)
	if err != nil {
		return err
	}
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", ecrData.AwsAuthType)

	log.Printf("[INFO] Created %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
	return nil
}

func castRawToAwsEcrAccessKeyData(data any) (ecrData api.AwsEcrAccessKeyData, err error) {
	dataJson, err := json.Marshal(data)
	err = json.Unmarshal(dataJson, &ecrData)
	return
}

func castRawToAwsEcrIamRoleData(data any) (ecrData api.AwsEcrIamRoleData, err error) {
	dataJson, err := json.Marshal(data)
	err = json.Unmarshal(dataJson, &ecrData)
	return
}

func resourceLaceworkIntegrationEcrCreateWithIAMRole(d *schema.ResourceData, lacework *api.Client) error {
	iamRoleData := api.AwsEcrIamRoleData{
		CrossAccountCredentials: api.AwsEcrCrossAccountCredentials{
			RoleArn:    d.Get("credentials.0.role_arn").(string),
			ExternalID: d.Get("credentials.0.external_id").(string),
		},
		LimitByTag:       castAttributeToStringSlice(d, "limit_by_tags"),
		LimitByRep:       castAttributeToStringSlice(d, "limit_by_repositories"),
		LimitNumImg:      d.Get("limit_num_imgs").(int),
		RegistryDomain:   d.Get("registry_domain").(string),
		NonOSPackageEval: d.Get("non_os_package_support").(bool),
	}

	labels := castAttributeToArrayKeyMapOfStrings(d, "limit_by_labels")
	if len(labels) != 0 {
		iamRoleData.LimitByLabel = labels
	}

	data := api.NewContainerRegistry(d.Get("name").(string),
		api.AwsEcrContainerRegistry,
		iamRoleData,
	)

	log.Printf("[INFO] Creating %s registry\n", api.AwsEcrContainerRegistry.String())
	response, err := lacework.V2.ContainerRegistries.Create(data)
	if err != nil {
		return err
	}

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
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

	ecrData, err := castRawToAwsEcrIamRoleData(response.Data.Data)
	if err != nil {
		return err
	}
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", ecrData.AwsAuthType)

	log.Printf("[INFO] Created %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
	return nil
}

func readEcrIam(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	response, err := lacework.V2.ContainerRegistries.GetAwsEcrIamRole(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	if response.Data.IntgGuid == d.Id() {
		d.Set("name", response.Data.Name)
		d.Set("intg_guid", response.Data.IntgGuid)
		d.Set("enabled", response.Data.Enabled == 1)
		d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
		d.Set("type_name", response.Data.Type)
		d.Set("org_level", response.Data.IsOrg == 1)

		d.Set("non_os_package_support", response.Data.Data.NonOSPackageEval)
		d.Set("registry_domain", response.Data.Data.RegistryDomain)
		d.Set("aws_auth_type", response.Data.Data.AwsAuthType)

		creds := make(map[string]string)
		creds["role_arn"] = response.Data.Data.CrossAccountCredentials.RoleArn
		creds["external_id"] = response.Data.Data.CrossAccountCredentials.ExternalID
		d.Set("credentials", []map[string]string{creds})

		d.Set("limit_num_imgs", response.Data.Data.LimitNumImg)
		if len(response.Data.Data.LimitByTag) != 0 {
			d.Set("limit_by_tags", response.Data.Data.LimitByTag)
		}
		if len(response.Data.Data.LimitByRep) != 0 {
			d.Set("limit_by_repositories", response.Data.Data.LimitByRep)
		}

		if len(response.Data.Data.LimitByLabel) != 0 {
			d.Set("limit_by_labels", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))
		}

		log.Printf("[INFO] Read %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), response.Data.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}

func readEcrAccessKey(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	response, err := lacework.V2.ContainerRegistries.GetAwsEcrAccessKey(d.Id())
	if err != nil {
		return resourceNotFound(d, err)
	}

	if response.Data.IntgGuid == d.Id() {
		d.Set("name", response.Data.Name)
		d.Set("intg_guid", response.Data.IntgGuid)
		d.Set("enabled", response.Data.Enabled == 1)
		d.Set("created_or_updated_time", response.Data.CreatedOrUpdatedTime)
		d.Set("created_or_updated_by", response.Data.CreatedOrUpdatedBy)
		d.Set("type_name", response.Data.Type)
		d.Set("org_level", response.Data.IsOrg == 1)

		d.Set("non_os_package_support", response.Data.Data.NonOSPackageEval)
		d.Set("registry_domain", response.Data.Data.RegistryDomain)
		d.Set("aws_auth_type", response.Data.Data.AwsAuthType)

		creds := make(map[string]string)
		creds["access_key_id"] = response.Data.Data.AccessKeyCredentials.AccessKeyID
		d.Set("credentials", []map[string]string{creds})

		d.Set("limit_num_imgs", response.Data.Data.LimitNumImg)
		if len(response.Data.Data.LimitByTag) != 0 {
			d.Set("limit_by_tags", response.Data.Data.LimitByTag)
		}
		if len(response.Data.Data.LimitByRep) != 0 {
			d.Set("limit_by_repositories", response.Data.Data.LimitByRep)
		}

		if len(response.Data.Data.LimitByLabel) != 0 {
			d.Set("limit_by_labels", castArrayOfStringKeyMapOfStringsToLimitByLabelSet(response.Data.Data.LimitByLabel))
		}

		log.Printf("[INFO] Read %s registry type with guid: %v\n", api.AwsEcrContainerRegistry.String(), response.Data.IntgGuid)
		return nil
	}

	d.SetId("")
	return nil
}
