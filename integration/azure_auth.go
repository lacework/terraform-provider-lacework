package integration

import (
	"encoding/json"
	"os"
)

type azureCredentialsFile struct {
	ClientSecret string `json:"keyId"`
	ClientID     string `json:"secret"`
}

func azureLoadDefaultCredentials() (azureCredentialsFile, error) {
	return azureLoadCredentials("AZURE_CREDENTIALS")
}

func azureLoadCredentials(envVar string) (c azureCredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
