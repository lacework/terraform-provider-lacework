package integration

import (
	"encoding/json"
	"os"
)

type googleCredentialsFile struct {
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientID     string `json:"client_id"`
	ProjectID    string `json:"project_id"`
}

func googleLoadDefaultCredentials() (googleCredentialsFile, error) {
	return googleLoadCredentials("GOOGLE_CREDENTIALS")
}

func googleLoadCredentials(envVar string) (c googleCredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
