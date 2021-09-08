package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func importLaceworkResourceGroup(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.ResourceGroupResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Resource Group with guid: %s\n", d.Id())

	if err := lacework.V2.ResourceGroups.Get(d.Id(), &response); err != nil {
		return nil, fmt.Errorf(
			"Unable to import Lacework resource. Resource Group with guid '%s' was not found.",
			d.Id(),
		)
	}
	log.Printf("[INFO] Resource Group  found using APIv2 with guid: %v\n", response.Data.ResourceGuid)
	return []*schema.ResourceData{d}, nil
}
