package lacework

import (
	"fmt"

	"github.com/lacework/go-sdk/v2/api"
)

// validateAgentlessScanningQueryText validates the LQL filter query used by
// Limited Workload Scanning before it is sent to the Lacework cloud account
// integration API.
//
// The integration API itself does not validate query_text, so without this
// check a syntactically invalid query (e.g. unbalanced braces, or a
// nonexistent LQL function like is_null()) is silently persisted, and the
// scanner falls back to scanning everything with no error surfaced anywhere.
//
// An empty query_text means "scan everything" and is intentionally not
// validated.
func validateAgentlessScanningQueryText(lacework *api.Client, queryText string) error {
	if queryText == "" {
		return nil
	}

	if _, err := lacework.V2.Query.Validate(api.ValidateQuery{QueryText: queryText}); err != nil {
		return fmt.Errorf("invalid query_text: %s", err)
	}
	return nil
}
