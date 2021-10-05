package integration

import (
	"os"
)

func datadogEnvVarsDefault() string {
	return os.Getenv("DATADOG_API_KEY")
}
