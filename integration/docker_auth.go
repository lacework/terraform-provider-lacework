package integration

import (
	"encoding/json"
	"os"
)

type dockerCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func dockerLoadDefaultCredentials() (dockerCredentials, error) {
	return dockerLoadCredentials("DOCKER_CREDENTIALS")
}

func dockerLoadCredentials(envVar string) (c dockerCredentials, err error) {
	err = json.Unmarshal([]byte(os.Getenv(envVar)), &c)
	return
}
