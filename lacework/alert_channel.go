package lacework

import (
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

// VerifyAlertChannel tests the integration of an alert channel
func VerifyAlertChannel(id string, lacework *api.Client) error {
	if err := lacework.V2.AlertChannels.Test(id); err != nil {
		// rollback terraform create upon error testing integration
		if _, deleteErr := lacework.Integrations.Delete(id); deleteErr != nil {
			return errors.Wrapf(err, "Unable to rollback changes: %v", deleteErr)
		}
		return err
	}
	return nil
}
