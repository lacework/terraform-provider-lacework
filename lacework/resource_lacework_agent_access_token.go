package lacework

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lacework/go-sdk/v2/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLaceworkAgentAccessToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAgentAccessTokenCreate,
		Read:   resourceLaceworkAgentAccessTokenRead,
		Update: resourceLaceworkAgentAccessTokenUpdate,
		Delete: resourceLaceworkAgentAccessTokenDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkAgentAccessToken,
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
	response, err := lacework.V2.AgentAccessTokens.Create(tokenName, tokenDesc)
	if err != nil {
		return err
	}

	token := response.Data
	d.SetId(token.TokenAlias)
	d.Set("name", token.TokenAlias)
	d.Set("token", token.AccessToken)
	d.Set("description", token.Props.Description)
	d.Set("version", token.Version)
	d.Set("enabled", token.State())
	d.Set("last_updated_time", token.CreatedTime.Format(time.RFC3339))
	d.Set("created_time", token.Props.CreatedTime.Format(time.RFC3339))

	// very unusual but, if the user creates a token disabled, update its status
	if !tokenEnabled {
		log.Println("[INFO] Disabling agent access token.")
		_, err = lacework.V2.AgentAccessTokens.Update(token.AccessToken, api.AgentAccessTokenRequest{Enabled: 0})
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
	response, err := lacework.V2.AgentAccessTokens.Get(d.Get("token").(string))
	if err != nil {
		return resourceNotFound(d, err)
	}

	token := response.Data
	if token.TokenAlias == d.Id() {
		d.Set("name", token.TokenAlias)
		d.Set("token", token.AccessToken)
		d.Set("description", token.Props.Description)
		d.Set("enabled", token.State())
		d.Set("version", token.Version)
		d.Set("last_updated_time", token.CreatedTime.Format(time.RFC3339))
		d.Set("created_time", token.Props.CreatedTime.Format(time.RFC3339))

		log.Printf("[INFO] Read agent access token. name=%s, description=%s, enabled=%t",
			token.TokenAlias, token.Props.Description, token.State())
		return nil
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAgentAccessTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		token    = api.AgentAccessTokenRequest{
			TokenAlias: d.Get("name").(string),
			Enabled:    0,
			Props: &api.AgentAccessTokenProps{
				Description: d.Get("description").(string),
			},
		}
	)

	if d.Get("enabled").(bool) {
		token.Enabled = 1
	}

	log.Printf("[INFO] Updating agent access token. name=%s, description=%s, enabled=%t",
		token.TokenAlias, token.Props.Description, d.Get("enabled").(bool))
	response, err := lacework.V2.AgentAccessTokens.Update(d.Get("token").(string), token)
	if err != nil {
		return err
	}

	nToken := response.Data
	d.SetId(token.TokenAlias)
	d.Set("name", nToken.TokenAlias)
	d.Set("token", nToken.AccessToken)
	d.Set("description", nToken.Props.Description)
	d.Set("enabled", nToken.State())
	d.Set("version", nToken.Version)
	d.Set("last_updated_time", nToken.CreatedTime.Format(time.RFC3339))
	d.Set("created_time", nToken.Props.CreatedTime.Format(time.RFC3339))

	log.Printf("[INFO] Agent access token updated")
	return nil
}

func resourceLaceworkAgentAccessTokenDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework  = meta.(*api.Client)
		tokenName = fmt.Sprintf("%s-%s-deleted", d.Get("name").(string), randomString(5))
	)

	// @afiune agent access tokens, by design, cannot be deleted, instead of deleting
	// them, we only disable them, but we will also modify its TokenAlias since that
	// field has a unique constraint. There can't be two tokens with the same alias.

	log.Printf("[INFO] Disabling agent access token. name=%s", tokenName)
	_, err := lacework.V2.AgentAccessTokens.Update(d.Get("token").(string), api.AgentAccessTokenRequest{Enabled: 0, TokenAlias: tokenName})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Agent access token disabled and updated with name '%s'.", tokenName)
	return nil
}

func importLaceworkAgentAccessToken(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing agent access token.")
	response, err := lacework.V2.AgentAccessTokens.Get(d.Id())
	if err != nil {
		return nil, err
	}

	if response.Data.AccessToken == d.Id() {
		log.Printf("[INFO] agent access token found. name=%s, description=%s, enabled=%t",
			response.Data.TokenAlias, response.Data.Props.Description, response.Data.State())

		d.Set("token", response.Data.AccessToken)
		d.SetId(response.Data.TokenAlias)

		return []*schema.ResourceData{d}, nil
	}

	log.Printf("[INFO] Raw response: %v\n", response)
	return nil, fmt.Errorf(
		"Unable to import Lacework resource. Agent access token '%s' was not found.",
		d.Id(),
	)
}
