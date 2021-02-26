package lacework

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAgentAccessToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAgentAccessTokenCreate,
		Read:   resourceLaceworkAgentAccessTokenRead,
		Update: resourceLaceworkAgentAccessTokenUpdate,
		Delete: resourceLaceworkAgentAccessTokenDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkAgentAccessToken,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceLaceworkAgentAccessTokenCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework     = meta.(*api.Client)
		tokenName    = d.Get("name").(string)
		tokenDesc    = d.Get("description").(string)
		tokenEnabled = d.Get("enabled").(bool)
	)

	log.Printf("[INFO] Creating agent access token. name=%s, description=%s, enabled=%t",
		tokenName, tokenDesc, tokenEnabled)
	response, err := lacework.Agents.CreateToken(tokenName, tokenDesc)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAgentTokenResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point in time, we know the data field has a value
	token := response.Data[0]
	d.SetId(token.TokenAlias)
	d.Set("name", token.TokenAlias)
	d.Set("token", token.AccessToken)
	d.Set("description", token.Props.Description)
	d.Set("account", token.Account)
	d.Set("version", token.Version)
	d.Set("enabled", token.Status())
	d.Set("last_updated_time", token.LastUpdatedTime.Format(time.RFC3339))
	d.Set("created_time", token.Props.CreatedTime.Format(time.RFC3339))

	// very unusual but, if the user creates a token disabled, update its status
	if !tokenEnabled {
		log.Println("[INFO] Disabling agent access token.")
		_, err = lacework.Agents.UpdateTokenStatus(token.AccessToken, false)
		if err != nil {
			return err
		}
		d.Set("enabled", false)
	}

	log.Printf("[INFO] Agent access token created.")
	return nil
}

func resourceLaceworkAgentAccessTokenRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading agent access token.")
	response, err := lacework.Agents.GetToken(d.Get("token").(string))
	if err != nil {
		return err
	}

	for _, token := range response.Data {
		if token.TokenAlias == d.Id() {
			d.Set("name", token.TokenAlias)
			d.Set("token", token.AccessToken)
			d.Set("description", token.Props.Description)
			d.Set("enabled", token.Status())
			d.Set("account", token.Account)
			d.Set("version", token.Version)
			d.Set("last_updated_time", token.LastUpdatedTime.Format(time.RFC3339))
			d.Set("created_time", token.Props.CreatedTime.Format(time.RFC3339))

			log.Printf("[INFO] Read agent access token. name=%s, description=%s, enabled=%t",
				token.TokenAlias, token.Props.Description, token.Status())
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAgentAccessTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		token    = api.AgentTokenRequest{
			TokenAlias: d.Get("name").(string),
			Enabled:    0,
			Props: &api.AgentTokenProps{
				Description: d.Get("description").(string),
			},
		}
	)

	if d.Get("enabled").(bool) {
		token.Enabled = 1
	}

	log.Printf("[INFO] Updating agent access token. name=%s, description=%s, enabled=%t",
		token.TokenAlias, token.Props.Description, d.Get("enabled").(bool))
	response, err := lacework.Agents.UpdateToken(d.Get("token").(string), token)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateAgentTokenResponse(&response)
	if err != nil {
		return err
	}

	// @afiune at this point in time, we know the data field has a value
	nToken := response.Data[0]
	d.SetId(token.TokenAlias)
	d.Set("name", nToken.TokenAlias)
	d.Set("token", nToken.AccessToken)
	d.Set("description", nToken.Props.Description)
	d.Set("enabled", nToken.Status())
	d.Set("account", nToken.Account)
	d.Set("version", nToken.Version)
	d.Set("last_updated_time", nToken.LastUpdatedTime.Format(time.RFC3339))
	d.Set("created_time", nToken.Props.CreatedTime.Format(time.RFC3339))

	log.Printf("[INFO] Agent access token updated")
	return nil
}

func resourceLaceworkAgentAccessTokenDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework  = meta.(*api.Client)
		tokenName = fmt.Sprintf("%s-%s-deleted", d.Get("name").(string), randomString(5))
		token     = api.AgentTokenRequest{
			TokenAlias: tokenName,
			Enabled:    0,
		}
	)

	// @afiune agent access tokens, by design, cannot be deleted, instead of deleting
	// them, we only disable them, but we will also modify its TokenAlias since that
	// field has a unique constraint. There can't be two tokens with the same alias.

	log.Printf("[INFO] Disabling agent access token. name=%s", tokenName)
	_, err := lacework.Agents.UpdateToken(d.Get("token").(string), token)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Agent access token disabled and updated with name '%s'.", tokenName)
	return nil
}

func importLaceworkAgentAccessToken(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing agent access token.")
	response, err := lacework.Agents.GetToken(d.Id())
	if err != nil {
		return nil, err
	}

	for _, token := range response.Data {
		if token.AccessToken == d.Id() {
			log.Printf("[INFO] agent access token found. name=%s, description=%s, enabled=%t",
				token.TokenAlias, token.Props.Description, token.Status())

			d.Set("token", token.AccessToken)
			d.SetId(token.TokenAlias)

			return []*schema.ResourceData{d}, nil
		}
	}

	log.Printf("[INFO] Raw response: %v\n", response)
	return nil, fmt.Errorf(
		"Unable to import Lacework resource. Agent access token '%s' was not found.",
		d.Id(),
	)
}

// validateAgentTokenResponse checks weather or not the server response has
// any inconsistent data, it returns a friendly error message describing the
// problem and how to report it
func validateAgentTokenResponse(response *api.AgentTokensResponse) error {
	if len(response.Data) == 0 {
		// @afiune this edge case should never happen, if we land here it means that
		// something went wrong in the server side of things (Lacework API), so let
		// us inform that to our users
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your agent token was
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	return nil
}
