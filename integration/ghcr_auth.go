package integration

import (
	"encoding/json"
	"os"
)

type ghcrCredentials struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func ghcrLoadDefaultCredentials() (ghcrCredentials, error) {
	return ghcrLoadCredentials("GHCR_CREDENTIALS")
}

func ghcrLoadCredentials(envVar string) (c ghcrCredentials, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
