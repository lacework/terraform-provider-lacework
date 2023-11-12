package integration

import (
	"encoding/json"
	"os"
)

type azureCredentialsFile struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
}

func azureLoadDefaultCredentials() (azureCredentialsFile, error) {
	return azureLoadCredentials("AZURE_CREDENTIALS")
}

func azureLoadCredentials(envVar string) (c azureCredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
