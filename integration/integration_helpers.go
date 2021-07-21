package integration

import (
	"github.com/lacework/go-sdk/api"
	"log"
	"os"
	"strings"
)

func LwTestCLient() (lw *api.Client) {
	lw, err := api.NewClient(os.Getenv("CI_ACCOUNT"),
		api.WithApiKeys(os.Getenv("CI_API_KEY"), os.Getenv("CI_API_SECRET")))

	if err != nil {
		log.Fatalf("Failed to create new go-sdk client, %v", err)
	}

	return
}

func GetIntegrationName(result string) string {
	resultSplit := strings.Split(result, "[id=")
	id := strings.Split(resultSplit[1], "]")[0]

	res, err := LwTestCLient().Integrations.Get(id)
	if err != nil {
		log.Fatalf("Unable to find integration id: %v", id)
	}

	return res.Data[0].Name
}
