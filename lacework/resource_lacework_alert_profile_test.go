package lacework

import (
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestAlertProfileTemplateToArrayOfType(t *testing.T) {
	var (
		input        = []map[string]string{{"name": "alertName", "event_name": "alertEventName", "description": "alertDescription", "subject": "alertSubject"}}
		d            = resourceLaceworkAlertProfile()
		testResource = d.TestResourceData()
		alerts       []api.AlertTemplate
		actual       = []api.AlertTemplate{{"alertName", "alertEventName", "alertDescription", "alertSubject"}}
	)

	testResource.Set("alert", input)

	err := castSchemaSetToArrayOfAlertTemplate(testResource, "alert", &alerts)
	assert.NoError(t, err)
	assert.Equal(t, alerts, actual,
		"%s did not match expected value: %s", actual, alerts,
	)
}
