package lacework

import (
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

// VerifyAlertChannel tests the integration of an alert channel
func VerifyAlertChannel(id string, lacework *api.Client) error {
	if err := lacework.V2.AlertChannels.Test(id); err != nil {
		// rollback terraform create upon error testing integration
		if deleteErr := lacework.V2.AlertChannels.Delete(id); deleteErr != nil {
			return errors.Wrapf(deleteErr, "Unable to rollback changes: %v", err)
		}
		return err
	}
	return nil
}
