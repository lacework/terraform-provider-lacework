package lacework

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestRetryWithIntervalRetriesThenSucceeds(t *testing.T) {
	attempts := 0
	start := time.Now()
	err := retryWithInterval(context.Background(), time.Minute, 20*time.Millisecond, func() *retry.RetryError {
		attempts++
		if attempts < 3 {
			return retry.RetryableError(errors.New("not yet"))
		}
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
	// two waits between three attempts
	assert.GreaterOrEqual(t, time.Since(start), 40*time.Millisecond)
	// the SDK backoff would take at least 1.5s here; the fixed interval must not
	assert.Less(t, time.Since(start), 500*time.Millisecond)
}

func TestRetryWithIntervalReturnsLastNonRetryableError(t *testing.T) {
	attempts := 0
	err := retryWithInterval(context.Background(), time.Minute, time.Millisecond, func() *retry.RetryError {
		attempts++
		if attempts < 2 {
			return retry.RetryableError(errors.New("transient"))
		}
		return retry.NonRetryableError(errors.New("gave up"))
	})
	assert.EqualError(t, err, "gave up")
	assert.Equal(t, 2, attempts)
}

// The retries DiffSuppressFunc must not drop the attribute from the create
// diff, otherwise d.Get("retries") is 0 at create time and only one attempt
// is made. TestResourceDataRaw builds the data the same way apply does.
func TestAzureRetriesDefaultSurvivesCreateDiff(t *testing.T) {
	for name, res := range map[string]*schema.Resource{
		"azure_cfg":                resourceLaceworkIntegrationAzureCfg(),
		"azure_al":                 resourceLaceworkIntegrationAzureActivityLog(),
		"azure_ad_al":              resourceLaceworkIntegrationAzureAdAl(),
		"azure_agentless_scanning": resourceLaceworkIntegrationAzureAgentlessScanning(),
		"azure_dspm":               resourceLaceworkAzureDspm(),
	} {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Equal(t, 30, d.Get("retries").(int), name)
		d = schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"retries": 10})
		assert.Equal(t, 10, d.Get("retries").(int), name)
	}
}
