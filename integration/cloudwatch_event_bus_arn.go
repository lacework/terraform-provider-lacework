package integration

import (
	"os"
)

func cloudwatchEnvVarsDefault() string {
	return os.Getenv("CLOUDWATCH_EVENT_BUS_ARN")
}
