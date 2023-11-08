package integration

import (
	"log"
	"os"
	"regexp"

	"github.com/lacework/go-sdk/api"
)

var (
	LwClient    *api.Client
	LwOrgClient *api.Client
	LwApiToken  string
	tokenEnvVar map[string]string
)

func init() {
	LwClient = lwTestCLient()
	LwOrgClient = lwOrgTestClient()
	LwApiToken = generateUniqueApiToken()
	tokenEnvVar = map[string]string{
		"LW_API_TOKEN": LwApiToken,
	}
}

func generateUniqueApiToken() string {
	token, err := LwClient.GenerateToken()
	if err != nil {
		log.Fatalf("Failed to generate API token, %v", err)
	}

	return token.Token
}

func lwTestCLient() (lw *api.Client) {
	lw, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiV2(),
	)

	if err != nil {
		log.Fatalf("Failed to create new go-sdk client, %v", err)
	}
	return
}

func lwOrgTestClient() (lw *api.Client) {
	lw, err := api.NewClient(os.Getenv("LW_ACCOUNT"),
		api.WithApiKeys(os.Getenv("LW_API_KEY"), os.Getenv("LW_API_SECRET")),
		api.WithSubaccount(os.Getenv("LW_SUBACCOUNT")),
		api.WithApiV2(),
		api.WithOrgAccess(),
	)

	if err != nil {
		log.Fatalf("Failed to create new go-sdk client, %v", err)
	}
	return
}

func GetCloudAccountIntegrationName(result string) string {
	var res api.CloudAccountResponse
	id := GetIDFromTerraResults(result)

	err := LwClient.V2.CloudAccounts.Get(id, &res)
	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data.Name
}

func GetCloudOrgAccountIntegrationName(result string) string {
	var res api.CloudAccountResponse
	id := GetIDFromTerraResults(result)

	err := LwOrgClient.V2.CloudAccounts.Get(id, &res)
	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data.Name
}

func GetCloudAccountEksAuditLogData(result string) api.AwsEksAuditData {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.CloudAccounts.GetAwsEksAudit(id)
	if err != nil {
		log.Fatalf("Unable to find eks audit log id: %s\n Response: %v", id, response)
	}

	return response.Data.Data
}

func GetCloudAccountAgentlessScanningResponse(result string) api.AwsSidekickResponse {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.CloudAccounts.GetAwsSidekick(id)
	if err != nil {
		log.Fatalf("Unable to find aws sidekick id: %s\n Response: %v", id, response)
	}

	return response
}

func GetAwsAgentlessOrgScanningResponse(result string) api.AwsSidekickOrgResponse {
	id := GetIDFromTerraResults(result)

	response, err := LwOrgClient.V2.CloudAccounts.GetAwsSidekickOrg(id)
	if err != nil {
		log.Fatalf("Unable to find aws sidekick id: %s\n Response: %v", id, response)
	}

	return response
}

func GetCloudAccountGkeAuditLogData(result string) api.GcpGkeAuditData {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.CloudAccounts.GetGcpGkeAudit(id)
	if err != nil {
		log.Fatalf("Unable to find gke audit log id: %s\n Response: %v", id, response)
	}

	return response.Data.Data
}

func GetCloudAccountGcpPubSubAuditLogData(result string) api.GcpAlPubSubIntegrationResponse {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.CloudAccounts.GetGcpAlPubSub(id)
	if err != nil {
		log.Fatalf("Unable to find gke audit log id: %s\n Response: %v", id, response)
	}

	return response
}

func GetAlertChannelName(result string) string {
	var res api.AlertChannelResponse

	id := GetIDFromTerraResults(result)

	err := LwClient.V2.AlertChannels.Get(id, &res)
	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data.Name
}

func GetEcrWithCrossAccountCreds(result string) api.AwsEcrIamRoleIntegration {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetAwsEcrIamRole(id)
	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data
}

func GetGcpCfgIntegration(result string) api.GcpCfgIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.CloudAccounts.GetGcpCfg(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryGar(result string) api.GcpGarIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetGcpGar(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetGcpAgentlessScanningResponse(result string) api.GcpSidekickIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.CloudAccounts.GetGcpSidekick(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetAzureAgentlessScanningResponse(result string) api.AzureSidekickIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.CloudAccounts.GetAzureSidekick(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryGcr(result string) api.GcpGcrIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetGcpGcr(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryDockerhub(result string) api.DockerhubIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetDockerhub(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryGhcr(result string) api.GhcrIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetGhcr(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryInlineScanner(result string) api.InlineScannerIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetInlineScanner(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetContainerRegisteryProxyScanner(result string) api.ProxyScannerIntegrationResponse {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.V2.ContainerRegistries.GetProxyScanner(id)

	if err != nil {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res
}

func GetResourceGroupV2Description(result string) string {
	id := GetIDFromTerraResults(result)

	var response api.ResourceGroupResponse
	err := LwClient.V2.ResourceGroups.Get(id, &response)

	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Description
}

func GetResourceGroupDescription(result string) string {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.ResourceGroups.GetAws(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props.Description
}

func GetAzureResourceGroupProps(result string) api.AzureResourceGroupProps {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.ResourceGroups.GetAzure(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetGcpResourceGroupProps(result string) api.GcpResourceGroupProps {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.ResourceGroups.GetGcp(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetContainerResourceGroupProps(result string) api.ContainerResourceGroupProps {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.ResourceGroups.GetContainer(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetTeamMember(result string) api.TeamMember {
	id := GetIDFromTerraResults(result)

	var res api.TeamMemberResponse
	err := LwClient.V2.TeamMembers.Get(id, &res)
	if err != nil {
		log.Fatalf("Unable to find team member id: %s\n Response: %v", id, res)
	}

	return res.Data
}

func GetOrgTeamMember(result string) api.TeamMember {
	id := GetIDFromTerraResults(result)

	var res api.TeamMemberResponse
	err := LwOrgClient.V2.TeamMembers.Get(id, &res)
	if err != nil {
		log.Fatalf("Unable to find team member id: %s\n Response: %v", id, res)
	}

	return res.Data
}

func GetMachineResourceGroupProps(result string) api.MachineResourceGroupProps {
	id := GetIDFromTerraResults(result)

	response, err := LwClient.V2.ResourceGroups.GetMachine(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetLwAccountResourceGroupProps(result string) api.LwAccountResourceGroupProps {
	id := GetIDFromTerraResults(result)

	response, err := LwOrgClient.V2.ResourceGroups.GetLwAccount(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetAlertChannelProps(result string) api.AlertChannelResponse {
	id := GetIDFromTerraResults(result)

	var data api.AlertChannelResponse
	err := LwClient.V2.AlertChannels.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve alert channel with id: %s", id)
	}
	return data
}

func GetVulnerabilityExceptionProps(result string) api.VulnerabilityExceptionResponse {
	id := GetSpecificIDFromTerraResults(1, result)
	var data api.VulnerabilityExceptionResponse
	err := LwClient.V2.VulnerabilityExceptions.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve vulnerability exception with id: %s", id)
	}
	return data
}

func GetAlertRuleProps(result string) api.AlertRuleResponse {
	id := GetSpecificIDFromTerraResults(2, result)

	var data api.AlertRuleResponse
	err := LwClient.V2.AlertRules.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve alert rule with id: %s", id)
	}
	return data
}

func GetAlertProfileProps(result string) api.AlertProfileResponse {
	id := GetSpecificIDFromTerraResults(1, result)

	var data api.AlertProfileResponse
	err := LwClient.V2.Alert.Profiles.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve alert profile with id: %s", id)
	}
	return data
}

func GetReportRuleProps(result string) api.ReportRuleResponse {
	id := GetSpecificIDFromTerraResults(3, result)

	var data api.ReportRuleResponse
	err := LwClient.V2.ReportRules.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve report rule with id: %s", id)
	}
	return data
}

// GetSpecificIDFromTerraResults returns the specific index id found in the Terraform output
func GetSpecificIDFromTerraResults(i int, result string) string {
	re := regexp.MustCompile(`\[id=(.*?)\]`)
	match := re.FindAllStringSubmatch(result, -1)
	if len(match) >= i {
		return match[i-1][1]
	}
	return ""
}

// GetIDFromTerraResults returns the first id found in the Terraform output
func GetIDFromTerraResults(result string) string {
	return GetSpecificIDFromTerraResults(1, result)
}

func GetQueryProps(result string) api.QueryResponse {
	id := GetSpecificIDFromTerraResults(1, result)

	resp, err := LwClient.V2.Query.Get(id)
	if err != nil {
		log.Fatalf("Unable to retrieve vulnerability exception with id: %s", id)
	}
	return resp
}

func GetPolicyProps(result string) api.PolicyResponse {
	id := GetSpecificIDFromTerraResults(1, result)

	resp, err := LwClient.V2.Policy.Get(id)
	if err != nil {
		log.Fatalf("Unable to retrieve policy with id: %s", id)
	}
	return resp
}

func GetPolicyPropsById(id string) api.PolicyResponse {
	resp, err := LwClient.V2.Policy.Get(id)
	if err != nil {
		log.Fatalf("Unable to retrieve policy with id: %s", id)
	}
	return resp
}

func GetPolicyExceptionProps(result string, policyID string) (resp api.PolicyExceptionResponse) {
	id := GetSpecificIDFromTerraResults(1, result)
	err := LwClient.V2.Policy.Exceptions.Get(policyID, id, &resp)
	if err != nil {
		log.Fatalf("Unable to retrieve policy exception with id: %s", id)
	}
	return
}

func GetDataExportRuleProps(result string) api.DataExportRuleResponse {
	id := GetSpecificIDFromTerraResults(1, result)

	data, err := LwClient.V2.DataExportRules.Get(id)
	if err != nil {
		log.Fatalf("Unable to retrieve data export rule with id: %s", id)
	}
	return data
}
