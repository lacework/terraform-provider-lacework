package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

func importLaceworkECRIntegration(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework integration with guid: %s\n", d.Id())
	response, err := lacework.Integrations.Get(d.Id())
	if err != nil {
		return nil, err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			log.Printf("[INFO] Integration found with guid: %v\n", integration.IntgGuid)

			// @afiune this field is important since it tells the resource which API functions to use
			// we will extract it from the raw integration response so that the READ functions populate
			// the rest of the fiels
			if awsAuthType, ok := integration.Data["AWS_AUTH_TYPE"]; ok {
				d.Set("aws_auth_type", awsAuthType.(string))
			} else {
				// if this field does not exist, the user might be importing a wrong integration type
				// (or the SCHEMA changed again without notice...)
				return nil, errors.New("AWS_AUTH_TYPE not found. Are you importing an ECR integration?")
			}

			log.Printf("[INFO] ECR integration found with guid: %v\n", integration.IntgGuid)
			return []*schema.ResourceData{d}, nil
		}
	}

	log.Printf("[INFO] Raw integration response: %v\n", response)
	return nil, fmt.Errorf(
		"Unable to import Lacework resource. Integration with guid '%s' was not found.",
		d.Id(),
	)
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
				Default:     false,
				Description: "Enable program language scanning",
			},

			// TODO @afiune remove these resources when we release v1.0
			"limit_by_tag": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "*",
				Description:   "A comma-separated list of image tags to limit the assessment of images with matching tags",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_tags` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_tags"},
			},
			"limit_by_label": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "*",
				Description:   "A comma-separated list of image labels to limit the assessment of images with matching labels",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_labels` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_labels"},
			},

			"limit_by_repos": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "A comma-separated list of repositories to assess",
				Deprecated:    "This attribute will be replaced by a new attribute `limit_by_repositories` in version 1.0 of the Lacework provider",
				ConflictsWith: []string{"limit_by_repositories"},
			},
			// END TODO @afiune

			"limit_by_tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A list of image tags to limit the assessment of images with matching tags",
				ConflictsWith: []string{"limit_by_tag"},
			},
			"limit_by_labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A key based map of labels to limit the assessment of images with matching key:value labels",
				ConflictsWith: []string{"limit_by_label"},
			},
			"limit_by_repositories": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
				Optional:      true,
				Description:   "A list of repositories to assess",
				ConflictsWith: []string{"limit_by_repos"},
			},
			"limit_num_imgs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The maximum number of newest container images to assess per repository",
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
				Description: "Wheter or not this integration is configured at the Organization level",
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
	lacework := meta.(*api.Client)

	switch d.Get("aws_auth_type").(string) {
	case api.AwsEcrAccessKey.String():
		log.Printf("[INFO] %s Authentication: %s\n", api.EcrRegistry.String(), api.AwsEcrAccessKey.String())
		return resourceLaceworkIntegrationEcrReadWithAccessKey(d, lacework)

	case api.AwsEcrIAM.String():
		log.Printf("[INFO] %s Authentication: %s\n", api.EcrRegistry.String(), api.AwsEcrIAM.String())
		return resourceLaceworkIntegrationEcrReadWithIAMRole(d, lacework)

	default:
		return errors.Errorf("Unsupported ECR authentication type '%s'.", d.Get("aws_auth_type").(string))
	}
}

func resourceLaceworkIntegrationEcrReadWithIAMRole(d *schema.ResourceData, lacework *api.Client) error {
	log.Printf("[INFO] Reading %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())
	response, err := lacework.Integrations.GetAwsEcrWithCrossAccount(d.Id())
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
			d.Set("non_os_package_support", integration.Data.NonOSPackageEval)
			// @afiune this field is important for updates since it will force a new resource
			d.Set("aws_auth_type", integration.Data.AwsAuthType)
			d.Set("registry_domain", integration.Data.RegistryDomain)

			creds := make(map[string]string)
			creds["role_arn"] = integration.Data.Credentials.RoleArn
			creds["external_id"] = integration.Data.Credentials.ExternalID
			d.Set("credentials", []map[string]string{creds})

			d.Set("limit_num_imgs", integration.Data.LimitNumImg)

			if _, ok := d.GetOk("limit_by_tags"); ok {
				d.Set("limit_by_tags", strings.Split(integration.Data.LimitByTag, ","))
			} else {
				d.Set("limit_by_tag", integration.Data.LimitByTag)
			}

			if _, ok := d.GetOk("limit_by_labels"); ok {
				d.Set("limit_by_labels", strings.Split(integration.Data.LimitByLabel, ","))
			} else {
				d.Set("limit_by_label", integration.Data.LimitByLabel)
			}

			if _, ok := d.GetOk("limit_by_repositories"); ok {
				d.Set("limit_by_repositories", strings.Split(integration.Data.LimitByRep, ","))
			} else {
				d.Set("limit_by_repos", integration.Data.LimitByRep)
			}

			log.Printf("[INFO] Read %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkIntegrationEcrReadWithAccessKey(d *schema.ResourceData, lacework *api.Client) error {
	log.Printf("[INFO] Reading %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())
	response, err := lacework.Integrations.GetAwsEcrWithAccessKey(d.Id())
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
			d.Set("non_os_package_support", integration.Data.NonOSPackageEval)

			// @afiune this field is important for updates since it will force a new resource
			d.Set("aws_auth_type", integration.Data.AwsAuthType)
			d.Set("registry_domain", integration.Data.RegistryDomain)

			creds := make(map[string]string)
			creds["access_key_id"] = integration.Data.Credentials.AccessKeyID
			d.Set("credentials", []map[string]string{creds})

			d.Set("limit_num_imgs", integration.Data.LimitNumImg)

			if _, ok := d.GetOk("limit_by_tags"); ok {
				d.Set("limit_by_tags", strings.Split(integration.Data.LimitByTag, ","))
			} else {
				d.Set("limit_by_tag", integration.Data.LimitByTag)
			}

			if _, ok := d.GetOk("limit_by_labels"); ok {
				d.Set("limit_by_labels", strings.Split(integration.Data.LimitByLabel, ","))
			} else {
				d.Set("limit_by_label", integration.Data.LimitByLabel)
			}

			if _, ok := d.GetOk("limit_by_repositories"); ok {
				d.Set("limit_by_repositories", strings.Split(integration.Data.LimitByRep, ","))
			} else {
				d.Set("limit_by_repos", integration.Data.LimitByRep)
			}

			log.Printf("[INFO] Read %s integration %s registry type with guid: %v\n",
				api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
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

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	limitByRepos := d.Get("limit_by_repos").(string)
	if repos := castAttributeToStringSlice(d, "limit_by_repositories"); len(repos) != 0 {
		limitByRepos = strings.Join(repos, ",")
	}

	data := api.NewAwsEcrWithCrossAccountIntegration(d.Get("name").(string),
		api.AwsEcrDataWithCrossAccountCreds{
			Credentials: api.AwsCrossAccountCreds{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
			AwsEcrCommonData: api.AwsEcrCommonData{
				LimitByTag:       limitByTags,
				LimitByLabel:     limitByLabels,
				LimitByRep:       limitByRepos,
				LimitNumImg:      d.Get("limit_num_imgs").(int),
				RegistryDomain:   d.Get("registry_domain").(string),
				NonOSPackageEval: d.Get("non_os_package_support").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration %s registry with authentication %s\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), data.Data.AwsAuthType)
	response, err := lacework.Integrations.UpdateAwsEcrWithCrossAccount(data)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Verifying server response")
	err = validateEcrWithIAMRoleIntegrationResponse(&response)
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
	d.Set("non_os_package_support", integration.Data.NonOSPackageEval)
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Updated %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())

	return nil
}

func resourceLaceworkIntegrationEcrUpdateWithAccessKey(d *schema.ResourceData, lacework *api.Client) error {

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	limitByRepos := d.Get("limit_by_repos").(string)
	if repos := castAttributeToStringSlice(d, "limit_by_repositories"); len(repos) != 0 {
		limitByRepos = strings.Join(repos, ",")
	}

	data := api.NewAwsEcrWithAccessKeyIntegration(d.Get("name").(string),
		api.AwsEcrDataWithAccessKeyCreds{
			Credentials: api.AwsEcrAccessKeyCreds{
				AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
				SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
			},
			AwsEcrCommonData: api.AwsEcrCommonData{
				LimitByTag:       limitByTags,
				LimitByLabel:     limitByLabels,
				LimitByRep:       limitByRepos,
				LimitNumImg:      d.Get("limit_num_imgs").(int),
				RegistryDomain:   d.Get("registry_domain").(string),
				NonOSPackageEval: d.Get("non_os_package_support").(bool),
			},
		},
	)

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	data.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration %s registry with authentication %s\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), data.Data.AwsAuthType)
	response, err := lacework.Integrations.UpdateAwsEcrWithAccessKey(data)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Verifying server response")
	err = validateEcrWithAccessKeyIntegrationResponse(&response)
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
	d.Set("non_os_package_support", integration.Data.NonOSPackageEval)
	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Updated %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())

	return nil
}

func resourceLaceworkIntegrationEcrDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())

	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), d.Id())

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
	data := api.NewAwsEcrWithAccessKeyIntegration(d.Get("name").(string),
		api.AwsEcrDataWithAccessKeyCreds{
			Credentials: api.AwsEcrAccessKeyCreds{
				AccessKeyID:     d.Get("credentials.0.access_key_id").(string),
				SecretAccessKey: d.Get("credentials.0.secret_access_key").(string),
			},
			AwsEcrCommonData: api.AwsEcrCommonData{
				LimitByTag:       d.Get("limit_by_tag").(string),
				LimitByLabel:     d.Get("limit_by_label").(string),
				LimitByRep:       d.Get("limit_by_repos").(string),
				LimitNumImg:      d.Get("limit_num_imgs").(int),
				RegistryDomain:   d.Get("registry_domain").(string),
				NonOSPackageEval: d.Get("non_os_package_support").(bool),
			},
		},
	)

	log.Printf("[INFO] Creating %s integration %s registry with authentication %s\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), data.Data.AwsAuthType)
	response, err := lacework.Integrations.CreateAwsEcrWithAccessKey(data)
	if err != nil {
		return err
	}

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Verifying server response")
	err = validateEcrWithAccessKeyIntegrationResponse(&response)
	if err != nil {
		return err
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
	d.Set("non_os_package_support", integration.Data.NonOSPackageEval)

	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Created %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
	return nil
}

func resourceLaceworkIntegrationEcrCreateWithIAMRole(d *schema.ResourceData, lacework *api.Client) error {

	limitByTags := d.Get("limit_by_tag").(string)
	if tags := castAttributeToStringSlice(d, "limit_by_tags"); len(tags) != 0 {
		limitByTags = strings.Join(tags, ",")
	}

	limitByLabels := d.Get("limit_by_label").(string)
	if labels := castAttributeToStringKeyMapOfStrings(d, "limit_by_labels"); len(labels) != 0 {
		limitByLabels = joinMapStrings(labels, ",")
	}

	limitByRepos := d.Get("limit_by_repos").(string)
	if repos := castAttributeToStringSlice(d, "limit_by_repositories"); len(repos) != 0 {
		limitByRepos = strings.Join(repos, ",")
	}

	data := api.NewAwsEcrWithCrossAccountIntegration(d.Get("name").(string),
		api.AwsEcrDataWithCrossAccountCreds{
			Credentials: api.AwsCrossAccountCreds{
				RoleArn:    d.Get("credentials.0.role_arn").(string),
				ExternalID: d.Get("credentials.0.external_id").(string),
			},
			AwsEcrCommonData: api.AwsEcrCommonData{
				LimitByTag:       limitByTags,
				LimitByLabel:     limitByLabels,
				LimitByRep:       limitByRepos,
				LimitNumImg:      d.Get("limit_num_imgs").(int),
				RegistryDomain:   d.Get("registry_domain").(string),
				NonOSPackageEval: d.Get("non_os_package_support").(bool),
			},
		},
	)

	log.Printf("[INFO] Creating %s integration %s registry with authentication %s\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), data.Data.AwsAuthType)
	response, err := lacework.Integrations.CreateAwsEcrWithCrossAccount(data)
	if err != nil {
		return err
	}

	if !d.Get("enabled").(bool) {
		data.Enabled = 0
	}

	log.Printf("[INFO] Verifying server response")
	err = validateEcrWithIAMRoleIntegrationResponse(&response)
	if err != nil {
		return err
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
	d.Set("non_os_package_support", integration.Data.NonOSPackageEval)

	// @afiune this field is important for updates since it will force a new resource
	d.Set("aws_auth_type", integration.Data.AwsAuthType)

	log.Printf("[INFO] Created %s integration %s registry type with guid: %v\n",
		api.ContainerRegistryIntegration.String(), api.EcrRegistry.String(), integration.IntgGuid)
	return nil
}

// validateEcrWithIAMRoleIntegrationResponse checks weather or not the server
// response has any inconsistent data, it returns a friendly error message describing
// the problem and how to report it
func validateEcrWithIAMRoleIntegrationResponse(
	response *api.AwsEcrWithCrossAccountIntegrationResponse) error {
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

// validateEcrWithAccessKeyIntegrationResponse checks weather or not the server
// response has any inconsistent data, it returns a friendly error message describing
// the problem and how to report it
func validateEcrWithAccessKeyIntegrationResponse(
	response *api.AwsEcrWithAccessKeyIntegrationResponse) error {
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
