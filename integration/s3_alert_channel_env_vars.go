package integration

import (
	"encoding/json"
	"os"
)

type s3CredentialsFile struct {
	RoleArn    string `json:"role_arn"`
	ExternalID string `json:"external_id"`
}

func s3LoadBucketArn() string {
	return os.Getenv("S3_BUCKET_ARN")
}

func s3LoadCredentials(envVar string) (s s3CredentialsFile, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &s)
	return
}
