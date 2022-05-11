package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func resourceLaceworkReportRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkReportRuleCreate,
		Read:   resourceLaceworkReportRuleRead,
		Update: resourceLaceworkReportRuleUpdate,
		Delete: resourceLaceworkReportRuleDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkReportRule,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the report rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the report rule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The state of the report rule",
			},
			"aws_compliance_reports": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pci": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cis_s3": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"soc": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"soc_rev2": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hipaa": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"iso_2700": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"nist_800_171_rev2": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"nist_800_53_rev4": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"gcp_compliance_reports": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pci": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"pci_rev2": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cis": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cis_12": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"soc": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"soc_rev2": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hipaa": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hipaa_rev2": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"iso_27001": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"k8s": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"azure_compliance_reports": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pci": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cis": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"soc": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"cis_131": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"daily_compliance_reports": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_cloudtrail": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"aws_compliance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"azure_activity_log": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"azure_compliance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"gcp_audit_trail": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"gcp_compliance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"openshift_compliance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"openshift_compliance_events": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"platform": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"host_security": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"weekly_snapshot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Weekly Snapshot Compliance Report type",
			},
			"severities": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Description: "List of severities for the report rule. Valid severities are:" +
					" Critical, High, Medium, Low, Info",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(cases.Title(language.English).String(strings.ToLower(val.(string))))
					},
					ValidateFunc: func(value interface{}, key string) ([]string, []error) {
						switch strings.ToLower(value.(string)) {
						case "critical", "high", "medium", "low", "info":
							return nil, nil
						default:
							return nil, []error{
								fmt.Errorf(
									"%s: can only be 'Critical', 'High', 'Medium', 'Low', 'Info'", key,
								),
							}
						}
					},
				},
			},
			"resource_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of resource groups for the report rule",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
			},
			"email_alert_channels": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "List of email alert channels for the report rule",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(val interface{}) string {
						return strings.TrimSpace(val.(string))
					},
				},
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkReportRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		resourceGroups = d.Get("resource_groups").(*schema.Set).List()
		severities     = api.NewReportRuleSeverities(castAttributeToStringSlice(d, "severities"))
		channels       = d.Get("email_alert_channels").(*schema.Set).List()
	)

	reportRule, err := api.NewReportRule(d.Get("name").(string),
		api.ReportRuleConfig{
			Description:        d.Get("description").(string),
			Severities:         severities,
			ResourceGroups:     castStringSlice(resourceGroups),
			EmailAlertChannels: castStringSlice(channels),
			NotificationTypes:  getReportRuleNotifications(d),
		},
	)

	if err != nil {
		return err
	}

	if !d.Get("enabled").(bool) {
		reportRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Creating report rule with data:\n%+v\n", reportRule)
	response, err := lacework.V2.ReportRules.Create(reportRule)
	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("guid", response.Data.Guid)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Created report rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkReportRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		response api.ReportRuleResponse
	)

	log.Printf("[INFO] Reading report rule with guid %s\n", d.Id())
	err := lacework.V2.ReportRules.Get(d.Id(), &response)
	if err != nil {
		return resourceNotFound(d, err)
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("guid", response.Data.Guid)
	d.Set("description", response.Data.Filter.Description)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type_name", response.Data.Type)
	d.Set("severities", api.NewReportRuleSeveritiesFromIntSlice(response.Data.Filter.Severity).ToStringSlice())
	d.Set("email_alert_channels", response.Data.EmailAlertChannels)
	d.Set("resource_groups", response.Data.Filter.ResourceGroups)
	setNotificationTypes(d, response.Data.ReportNotificationTypes)

	log.Printf("[INFO] Read report rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkReportRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		resourceGroups = d.Get("resource_groups").(*schema.Set).List()
		severities     = api.NewReportRuleSeverities(castAttributeToStringSlice(d, "severities"))
		channels       = d.Get("email_alert_channels").(*schema.Set).List()
	)

	reportRule, err := api.NewReportRule(d.Get("name").(string),
		api.ReportRuleConfig{
			Description:        d.Get("description").(string),
			Severities:         severities,
			ResourceGroups:     castStringSlice(resourceGroups),
			EmailAlertChannels: castStringSlice(channels),
			NotificationTypes:  getReportRuleNotifications(d),
		},
	)

	if err != nil {
		return err
	}

	reportRule.Guid = d.Id()

	if !d.Get("enabled").(bool) {
		reportRule.Filter.Enabled = 0
	}

	log.Printf("[INFO] Updating report rule with data:\n%+v\n", reportRule)
	response, err := lacework.V2.ReportRules.Update(reportRule)
	if err != nil {
		return err
	}

	d.SetId(response.Data.Guid)
	d.Set("name", response.Data.Filter.Name)
	d.Set("description", response.Data.Filter.Description)
	d.Set("guid", response.Data.Guid)
	d.Set("enabled", response.Data.Filter.Enabled == 1)
	d.Set("created_or_updated_time", response.Data.Filter.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", response.Data.Filter.CreatedOrUpdatedBy)
	d.Set("type", response.Data.Type)

	log.Printf("[INFO] Updated report rule with guid %s\n", response.Data.Guid)
	return nil
}

func resourceLaceworkReportRuleDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting report rule with guid %s\n", d.Id())
	err := lacework.V2.ReportRules.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted report rule with guid %s\n", d.Id())
	return nil
}

func importLaceworkReportRule(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.ReportRuleResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Report Rule with guid: %s\n", d.Id())

	if err := lacework.V2.ReportRules.Get(d.Id(), &response); err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Report Rule with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Report Rule found with guid: %s\n", response.Data.Guid)
	return []*schema.ResourceData{d}, nil
}

func getReportRuleNotifications(d *schema.ResourceData) api.ReportRuleNotifications {
	var reports api.ReportRuleNotifications

	if _, ok := d.GetOk("aws_compliance_reports"); ok {
		awsReport := api.AwsReportRuleNotifications{
			AwsCisS3:          d.Get("aws_compliance_reports.0.cis_s3").(bool),
			AwsHipaa:          d.Get("aws_compliance_reports.0.hipaa").(bool),
			AwsIso2700:        d.Get("aws_compliance_reports.0.iso_2700").(bool),
			AwsNist80053Rev4:  d.Get("aws_compliance_reports.0.nist_800_53_rev4").(bool),
			AwsNist800171Rev2: d.Get("aws_compliance_reports.0.nist_800_171_rev2").(bool),
			AwsPci:            d.Get("aws_compliance_reports.0.pci").(bool),
			AwsSoc:            d.Get("aws_compliance_reports.0.soc").(bool),
			AwsSocRev2:        d.Get("aws_compliance_reports.0.soc_rev2").(bool),
		}
		reports = append(reports, awsReport)
	}

	if _, ok := d.GetOk("gcp_compliance_reports"); ok {
		gcpReport := api.GcpReportRuleNotifications{
			GcpCis:       d.Get("gcp_compliance_reports.0.cis").(bool),
			GcpHipaa:     d.Get("gcp_compliance_reports.0.hipaa").(bool),
			GcpHipaaRev2: d.Get("gcp_compliance_reports.0.hipaa_rev2").(bool),
			GcpIso27001:  d.Get("gcp_compliance_reports.0.iso_27001").(bool),
			GcpCis12:     d.Get("gcp_compliance_reports.0.cis_12").(bool),
			GcpK8s:       d.Get("gcp_compliance_reports.0.k8s").(bool),
			GcpPci:       d.Get("gcp_compliance_reports.0.pci").(bool),
			GcpPciRev2:   d.Get("gcp_compliance_reports.0.pci_rev2").(bool),
			GcpSoc:       d.Get("gcp_compliance_reports.0.soc").(bool),
			GcpSocRev2:   d.Get("gcp_compliance_reports.0.soc_rev2").(bool),
		}
		reports = append(reports, gcpReport)
	}

	if _, ok := d.GetOk("azure_compliance_reports"); ok {
		azureReport := api.AzureReportRuleNotifications{
			AzureCis:    d.Get("azure_compliance_reports.0.cis").(bool),
			AzureCis131: d.Get("azure_compliance_reports.0.cis_131").(bool),
			AzurePci:    d.Get("azure_compliance_reports.0.pci").(bool),
			AzureSoc:    d.Get("azure_compliance_reports.0.soc").(bool),
		}
		reports = append(reports, azureReport)
	}

	if _, ok := d.GetOk("daily_compliance_reports"); ok {
		dailyReport := api.DailyEventsReportRuleNotifications{
			AgentEvents:               d.Get("daily_compliance_reports.0.host_security").(bool),
			OpenShiftCompliance:       d.Get("daily_compliance_reports.0.openshift_compliance").(bool),
			OpenShiftComplianceEvents: d.Get("daily_compliance_reports.0.openshift_compliance_events").(bool),
			PlatformEvents:            d.Get("daily_compliance_reports.0.host_security").(bool),
			AwsCloudtrailEvents:       d.Get("daily_compliance_reports.0.aws_cloudtrail").(bool),
			AwsComplianceEvents:       d.Get("daily_compliance_reports.0.host_security").(bool),
			AzureComplianceEvents:     d.Get("daily_compliance_reports.0.aws_compliance").(bool),
			AzureActivityLogEvents:    d.Get("daily_compliance_reports.0.azure_activity_log").(bool),
			GcpAuditTrailEvents:       d.Get("daily_compliance_reports.0.gcp_audit_trail").(bool),
			GcpComplianceEvents:       d.Get("daily_compliance_reports.0.gcp_compliance").(bool),
		}
		reports = append(reports, dailyReport)
	}

	if _, ok := d.GetOk("weekly_snapshot"); ok {
		weeklyReport := api.WeeklyEventsReportRuleNotifications{
			TrendReport: d.Get("weekly_snapshot").(bool),
		}
		reports = append(reports, weeklyReport)
	}

	return reports
}

func setNotificationTypes(d *schema.ResourceData, notifications api.ReportRuleNotificationTypes) {
	d.Set("aws_compliance_reports.0.cis_s3", notifications.AwsCisS3)
	d.Set("aws_compliance_reports.0.hipaa", notifications.AwsHipaa)
	d.Set("aws_compliance_reports.0.iso_2700", notifications.AwsIso2700)
	d.Set("aws_compliance_reports.0.nist_800_53_rev4", notifications.AwsNist80053Rev4)
	d.Set("aws_compliance_reports.0.nist_800_171_rev2", notifications.AwsNist800171Rev2)
	d.Set("aws_compliance_reports.0.pci", notifications.AwsPci)
	d.Set("aws_compliance_reports.0.soc", notifications.AwsSoc)
	d.Set("aws_compliance_reports.0.soc_rev2", notifications.AwsSocRev2)

	d.Set("gcp_compliance_reports.0.cis", notifications.GcpCis)
	d.Set("gcp_compliance_reports.0.hipaa", notifications.GcpHipaa)
	d.Set("gcp_compliance_reports.0.hipaa_rev2", notifications.GcpHipaaRev2)
	d.Set("gcp_compliance_reports.0.iso_27001", notifications.GcpIso27001)
	d.Set("gcp_compliance_reports.0.cis_12", notifications.GcpCis12)
	d.Set("gcp_compliance_reports.0.k8s", notifications.GcpK8s)
	d.Set("gcp_compliance_reports.0.pci", notifications.GcpPci)
	d.Set("gcp_compliance_reports.0.pci_rev2", notifications.GcpPciRev2)
	d.Set("gcp_compliance_reports.0.soc", notifications.GcpSoc)
	d.Set("gcp_compliance_reports.0.soc_rev2", notifications.GcpSocRev2)

	d.Set("azure_compliance_reports.0.cis", notifications.AzureCis)
	d.Set("azure_compliance_reports.0.cis_131", notifications.AzureCis131)
	d.Set("azure_compliance_reports.0.pci", notifications.AzurePci)
	d.Set("azure_compliance_reports.0.soc", notifications.AzureSoc)

	d.Set("daily_compliance_reports.0.host_security", notifications.AgentEvents)
	d.Set("daily_compliance_reports.0.openshift_compliance", notifications.OpenShiftCompliance)
	d.Set("daily_compliance_reports.0.openshift_compliance_events", notifications.OpenShiftComplianceEvents)
	d.Set("daily_compliance_reports.0.host_security", notifications.PlatformEvents)
	d.Set("daily_compliance_reports.0.aws_cloudtrail", notifications.AwsCloudtrailEvents)
	d.Set("daily_compliance_reports.0.host_security", notifications.AwsComplianceEvents)
	d.Set("daily_compliance_reports.0.aws_compliance", notifications.AzureComplianceEvents)
	d.Set("daily_compliance_reports.0.azure_activity_log", notifications.AzureActivityLogEvents)
	d.Set("daily_compliance_reports.0.gcp_audit_trail", notifications.GcpAuditTrailEvents)
	d.Set("daily_compliance_reports.0.gcp_compliance", notifications.GcpComplianceEvents)

	d.Set("weekly_snapshot", notifications.TrendReport)
}
