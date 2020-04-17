package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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

	log.Printf("[INFO] Raw integration response: %v\n", response)
	return nil, fmt.Errorf(
		"Unable to import Lacework resource. Integration with guid '%s' was not found.",
		d.Id(),
	)
}
