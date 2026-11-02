package lacework

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

// A telemetry outage (HTTP 503 from OtelMetrics) must not fail the read: it stopped the
// terraform destroy of a self-deployment rollback (CAD-2529).
func TestDataLaceworkMetricModuleReadIgnoresSendError(t *testing.T) {
	t.Setenv(api.DisableTelemetry, "") // Send is a no-op when this is set

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client, err := api.NewClient("test", api.WithURL(srv.URL), api.WithToken("token"))
	if err != nil {
		t.Fatal(err)
	}

	d := schema.TestResourceDataRaw(t, dataSourceLaceworkMetricModule().Schema,
		map[string]interface{}{"name": "terraform-azure-agentless-scanning", "version": "2.0.0"})

	if err := dataLaceworkMetricModuleRead(d, client); err != nil {
		t.Fatalf("expected no error on a telemetry failure, got: %v", err)
	}
	if hits.Load() == 0 {
		t.Fatal("the telemetry endpoint was never called")
	}
	if d.Id() == "" {
		t.Fatal("expected the ID to be set when the send fails")
	}
}
