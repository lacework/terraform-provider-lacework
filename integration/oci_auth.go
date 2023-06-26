package integration

import (
	"encoding/json"
	"os"
)

type ociCredentialsFile struct {
	User        string
	Fingerprint string
	TenacyID    string `json:"tenacy_id"`
	TenacyName  string `json:"tenacy_name"`
	Region      string
	PrivateKey  string `json:"private_key"`
}

func ociLoadDefaultCredentials() (ociCredentialsFile, error) {
	var c ociCredentialsFile
	err := json.Unmarshal([]byte(os.Getenv("OCI_CREDENTIALS")), &c)
	return c, err
}
