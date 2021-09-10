package integration

import (
	"log"
	"os"
	"strings"

	"github.com/lacework/go-sdk/api"
)

var LwClient *api.Client

func init() {
	LwClient = lwTestCLient()
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

func GetIntegrationName(result string) string {
	resultSplit := strings.Split(result, "[id=")
	id := strings.Split(resultSplit[1], "]")[0]

	res, err := LwClient.Integrations.Get(id)
	if err != nil || len(res.Data) == 0 {
		log.Fatalf("Unable to find integration id: %s\n Response: %v", id, res)
	}

	return res.Data[0].Name
}

func GetResourceGroupDescription(result string) string {
	resultSplit := strings.Split(result, "[id=")
	id := strings.Split(resultSplit[1], "]")[0]

	response, err := LwClient.V2.ResourceGroups.GetAws(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props.Description
}

func GetAzureResourceGroupProps(result string) api.AzureResourceGroupProps {
	resultSplit := strings.Split(result, "[id=")
	id := strings.Split(resultSplit[1], "]")[0]

	response, err := LwClient.V2.ResourceGroups.GetAzure(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}

func GetGcpResourceGroupProps(result string) api.GcpResourceGroupProps {
	resultSplit := strings.Split(result, "[id=")
	id := strings.Split(resultSplit[1], "]")[0]

	response, err := LwClient.V2.ResourceGroups.GetGcp(id)
	if err != nil {
		log.Fatalf("Unable to find resource group id: %s\n Response: %v", id, response)
	}

	return response.Data.Props
}
