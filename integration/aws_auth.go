package integration

import (
	"encoding/json"
	"os"
)

type awsCredentialsFile struct {
	RoleArn    string `json:"role_arn"`
	ExternalID string `json:"external_id"`
}

func awsLoadDefaultCredentials() (awsCredentialsFile, error) {
	return awsLoadCredentials("AWS_CREDS")
}

func awsLoadCredentials(envVar string) (c awsCredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
