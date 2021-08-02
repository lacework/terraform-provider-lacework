package lacework

import (
	"github.com/lacework/go-sdk/api"
)

// TestAlertChannel tests the integration of an alert channel
func TestAlertChannel(id string, meta interface{}) error {
	lacework := meta.(*api.Client)
	err := lacework.V2.AlertChannels.Test(id)

	// rollback terraform create upon error testing integration
	if err != nil {
		_, err := lacework.Integrations.Delete(id)
		return err
	}
	return nil
}
