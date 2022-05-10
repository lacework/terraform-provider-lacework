package lacework

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func notFound(err error) bool {
	var e *resource.NotFoundError
	return errors.As(err, &e)
}

func resourceNotFound(d *schema.ResourceData, err error, id string) error {
	if notFound(err) && !d.IsNewResource() {
		log.Printf("[WARN] resource with guid: %s\n not found, removing from state", id)
		d.SetId("")
		return nil
	}
	return err
}
