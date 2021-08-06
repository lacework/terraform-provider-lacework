package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func importLaceworkIntegration(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework integration with guid: %s\n", d.Id())
	response, err := lacework.Integrations.Get(d.Id())
	if err != nil {
		return nil, err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			log.Printf("[INFO] Integration found with guid: %v\n", integration.IntgGuid)
			return []*schema.ResourceData{d}, nil
		}
	}

	log.Printf("[INFO] Raw APIv1 integration response: %v\n", response)

	log.Println("[WARN] Trying APIv2")
	var cloudAccount api.CloudAccountRaw
	if err := lacework.V2.CloudAccounts.Get(d.Id(), &cloudAccount); err != nil {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Integration with guid '%s' was not found.",
			d.Id(),
		)
	}
	log.Printf("[INFO] Cloud account integration found using APIv2 with guid: %v\n", cloudAccount.IntgGuid)
	return []*schema.ResourceData{d}, nil

}
