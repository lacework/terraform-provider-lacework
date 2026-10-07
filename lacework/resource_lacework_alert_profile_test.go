package lacework

import (
	"testing"

	"github.com/lacework/go-sdk/v2/api"
	"github.com/stretchr/testify/assert"
)

func TestOwnAlertTemplates(t *testing.T) {
	parentDefault := api.AlertTemplate{Name: "DefaultAlert", EventName: "Default", Description: "d", Subject: "s"}
	parentViolation := api.AlertTemplate{Name: "Violation", EventName: "Violation", Description: "d", Subject: "s"}
	own := api.AlertTemplate{Name: "Custom", EventName: "Custom Event", Description: "cd", Subject: "cs"}
	overridden := api.AlertTemplate{Name: "Violation", EventName: "Violation", Description: "changed", Subject: "s"}
	parent := []api.AlertTemplate{parentDefault, parentViolation}

	// inherited alerts are dropped, own and overridden ones kept in schema shape
	got := ownAlertTemplates([]api.AlertTemplate{parentDefault, own, overridden}, parent, nil)
	assert.Equal(t, []map[string]interface{}{
		{"name": "Custom", "event_name": "Custom Event", "description": "cd", "subject": "cs"},
		{"name": "Violation", "event_name": "Violation", "description": "changed", "subject": "s"},
	}, got)

	// an alert identical to the parent's is kept when it is already in state
	got = ownAlertTemplates([]api.AlertTemplate{parentDefault, parentViolation}, parent, []api.AlertTemplate{parentViolation})
	assert.Equal(t, []map[string]interface{}{
		{"name": "Violation", "event_name": "Violation", "description": "d", "subject": "s"},
	}, got)
}
