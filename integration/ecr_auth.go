package integration

import (
	"encoding/json"
	"os"
)

type ecrCredentialsFile struct {
	RoleArn        string `json:"role_arn"`
	ExternalID     string `json:"external_id"`
	RegistryDomain string `json:"registry_domain"`
}

func ecrLoadDefaultCredentials() (ecrCredentialsFile, error) {
	return ecrLoadCredentials("AWS_ECR_IAM")
}

func ecrLoadCredentials(envVar string) (c ecrCredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
