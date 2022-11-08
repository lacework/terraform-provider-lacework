package lacework

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func importLaceworkContainerRegistry(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)
	var response api.ContainerRegistryResponse

	log.Printf("[INFO] Importing Lacework Container Registry with guid: %s\n", d.Id())
	err := lacework.V2.CloudAccounts.Get(d.Id(), &response)
	if err != nil {
		return nil, err
	}

	if response.Data.IntgGuid == d.Id() {
		log.Printf("[INFO] Container Registry found using APIv2 with guid: %v\n", response.Data.IntgGuid)
		return []*schema.ResourceData{d}, nil
	}

	log.Printf("[INFO] Raw APIv2 Container Registry response: %v\n", response)

	return []*schema.ResourceData{d}, nil
}
