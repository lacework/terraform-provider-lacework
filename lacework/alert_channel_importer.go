package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/v2/api"
)

// VerifyAlertChannelAndRollback will test the integration of an alert channel,
// if the test is not successful, it will remove the alert channel (rollback)
func VerifyAlertChannelAndRollback(d *schema.ResourceData, lacework *api.Client) error {
	if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
		defer d.SetId("")
		// rollback terraform create upon error testing integration
		if deleteErr := lacework.V2.AlertChannels.Delete(d.Id()); deleteErr != nil {
			return errors.Wrapf(deleteErr, "Unable to rollback changes: %v", err)
		}
		return err
	}
	return nil
}
