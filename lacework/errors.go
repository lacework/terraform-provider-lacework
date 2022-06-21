package lacework

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func notFound(err error) bool {
	return strings.Contains(err.Error(), "404")
}

func resourceNotFound(d *schema.ResourceData, err error) error {
	if notFound(err) && !d.IsNewResource() {
		log.Printf("[WARN] resource with guid: %s\n not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return err
}
