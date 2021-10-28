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
)

func init() {
	LwClient = lwTestCLient()
	LwOrgClient = lwOrgTestClient()
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

func GetIntegrationName(result string) string {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.Integrations.Get(id)
	if err != nil || len(res.Data) == 0 {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data[0].Name
}

func GetEcrWithCrossAccountCreds(result string) api.AwsEcrWithCrossAccountIntegration {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.Integrations.GetAwsEcrWithCrossAccount(id)
	if err != nil || len(res.Data) == 0 {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data[0]
}

func GetContainerRegistryIntegration(result string) api.ContainerRegIntegration {
	id := GetIDFromTerraResults(result)

	res, err := LwClient.Integrations.GetContainerRegistry(id)

	if err != nil || len(res.Data) == 0 {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data[0]
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

func GetAlertRuleProps(result string) api.AlertRuleResponse {
	id := GetIDFromTerraResults(result)

	var data api.AlertRuleResponse
	err := LwClient.V2.AlertRules.Get(id, &data)
	if err != nil {
		log.Fatalf("Unable to retrieve alert rule with id: %s", id)
	}
	return data
}

func GetIDFromTerraResults(result string) string {
	re := regexp.MustCompile("\\[id=(.*?)\\]")
	match := re.FindStringSubmatch(result)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}
