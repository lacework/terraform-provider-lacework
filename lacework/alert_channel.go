package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
)

// VerifyAlertChannelAndRollback will test the integration of an alert channel,
// if the test is not successful, it will remove the alert channel (rollback)
func VerifyAlertChannelAndRollback(id string, lacework *api.Client) error {
	if err := lacework.V2.AlertChannels.Test(id); err != nil {
		// rollback terraform create upon error testing integration
		if deleteErr := lacework.V2.AlertChannels.Delete(id); deleteErr != nil {
			return errors.Wrapf(deleteErr, "Unable to rollback changes: %v", err)
		}
		return err
	}
	return nil
}

func importLaceworkAlertChannel(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.AlertChannelResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Alert Channel with guid: %s\n", d.Id())

	if err := lacework.V2.AlertChannels.Get(d.Id(), &response); err != nil {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Alert Channel with guid '%s' was not found.",
			d.Id(),
		)
	}
	log.Printf("[INFO] Alert Channel found with guid: %v\n", response.Data.IntgGuid)
	return []*schema.ResourceData{d}, nil
}
