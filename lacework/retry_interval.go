package lacework

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// azureCreateRetryInterval is the wait between attempts to create an Azure
// integration. The Lacework API validates the integration's credentials
// against Azure Storage synchronously, and Azure documents up to 10 minutes
// for a new RBAC role assignment to take effect, so the retry window has to
// be measured in minutes rather than the seconds the SDK backoff gives.
const azureCreateRetryInterval = 20 * time.Second

// retryWithInterval behaves like retry.RetryContext but waits a fixed
// interval between attempts instead of the SDK's 0.5s to 10s backoff, so a
// modest retries count can span several minutes.
func retryWithInterval(ctx context.Context, timeout, interval time.Duration, f retry.RetryFunc) error {
	var (
		resultErr   error
		resultErrMu sync.Mutex
	)

	c := &retry.StateChangeConf{
		Pending:      []string{"retryableerror"},
		Target:       []string{"success"},
		Timeout:      timeout,
		PollInterval: interval,
		Refresh: func() (interface{}, string, error) {
			rerr := f()

			resultErrMu.Lock()
			defer resultErrMu.Unlock()

			if rerr == nil {
				resultErr = nil
				return 42, "success", nil
			}

			resultErr = rerr.Err

			if rerr.Retryable {
				return 42, "retryableerror", nil
			}
			return nil, "quit", rerr.Err
		},
	}

	_, waitErr := c.WaitForStateContext(ctx)

	resultErrMu.Lock()
	defer resultErrMu.Unlock()

	if resultErr == nil {
		return waitErr
	}
	return resultErr
}
