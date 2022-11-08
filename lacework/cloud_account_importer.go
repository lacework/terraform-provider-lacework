package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func importLaceworkCloudAccount(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)
	var response api.CloudAccountResponse

	log.Printf("[INFO] Importing Lacework Cloud Account with guid: %s\n", d.Id())
	err := lacework.V2.CloudAccounts.Get(d.Id(), &response)
	if err != nil {
		return nil, err
	}

	if response.Data.IntgGuid == d.Id() {
		log.Printf("[INFO] Cloud Account integration found using APIv2 with guid: %v\n", response.Data.IntgGuid)
		return []*schema.ResourceData{d}, nil
	}

	log.Printf("[INFO] Raw APIv2 Cloud Account response: %v\n", response)
	return []*schema.ResourceData{d}, nil
}
