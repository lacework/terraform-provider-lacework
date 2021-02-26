package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkAgentAccessToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLaceworkAgentAccessTokenRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceLaceworkAgentAccessTokenRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Lookup agent access token.")
	response, err := lacework.Agents.ListTokens()
	if err != nil {
		return err
	}

	lookupName := d.Get("name").(string)
	for _, token := range response.Data {
		if token.TokenAlias == lookupName {
			log.Printf("[INFO] agent access token found. name=%s, description=%s, enabled=%t",
				token.TokenAlias, token.Props.Description, token.Status())

			d.Set("token", token.AccessToken)
			d.SetId(token.TokenAlias)

			return nil
		}
	}

	return fmt.Errorf("Agent access token with name '%s' was not found.", lookupName)
}
