package lacework

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkApiToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLaceworkApiTokenRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLaceworkApiTokenRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	response, err := lacework.GenerateToken()
	if err != nil {
		// return the api client error directly since it is user friendly
		return err
	}

	d.SetId(time.Now().UTC().String())
	d.Set("token", response.Token())

	return nil
}
