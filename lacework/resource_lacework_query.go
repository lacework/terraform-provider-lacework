package lacework

import (
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

func resourceLaceworkQuery() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkQueryCreate,
		Read:   resourceLaceworkQueryRead,
		Update: resourceLaceworkQueryUpdate,
		Delete: resourceLaceworkQueryDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkQuery,
		},

		Schema: map[string]*schema.Schema{
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The query string",
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"result_schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getIDFromLQLQuery(query string) (string, error) {
	if strings.Contains(query, "{") {
		rexQueryID := regexp.MustCompile(`^([^{]*)`)
		if id := rexQueryID.FindStringSubmatch(query); id != nil {
			queryID := strings.TrimSpace(id[0])
			if queryID != "" && !strings.Contains(queryID, " ") {
				return queryID, nil
			}
		}
	}

	return "", errors.New(`query id not found. (malformed)

> Your query:
` + query + `

> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`)
}

func resourceLaceworkQueryCreate(d *schema.ResourceData, meta interface{}) error {
	var lacework = meta.(*api.Client)

	queryID, err := getIDFromLQLQuery(d.Get("query").(string))
	if err != nil {
		return err
	}

	query := api.NewQuery{
		QueryID:   queryID,
		QueryText: d.Get("query").(string),
	}

	log.Printf("[INFO] Creating Query with data:\n%+v\n", query)
	response, err := lacework.V2.Query.Create(query)
	if err != nil {
		return err
	}

	d.SetId(response.Data.QueryID)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("result_schema", response.Data.ResultSchema)

	log.Printf("[INFO] Created Query with guid %s\n", response.Data.QueryID)
	return nil
}

func resourceLaceworkQueryRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	log.Printf("[INFO] Reading Query with guid %s\n", d.Id())
	response, err := lacework.V2.Query.Get(d.Id())
	if err != nil {
		return err
	}

	d.Set("query", response.Data.QueryText)
	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("result_schema", response.Data.ResultSchema)

	log.Printf("[INFO] Read Query with guid %s\n", response.Data.QueryID)
	return nil
}

func resourceLaceworkQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	queryID, err := getIDFromLQLQuery(d.Get("query").(string))
	if err != nil {
		return err
	}

	if d.Id() != queryID {
		return errors.Errorf(
			"unable to change id of an existing query.\n\nOld ID: %s\n\n New ID: %s",
			d.Id(), queryID,
		)
	}

	query := api.UpdateQuery{
		QueryText: d.Get("query").(string),
	}

	log.Printf("[INFO] Updating Query with data:\n%+v\n", query)
	response, err := lacework.V2.Query.Update(d.Id(), query)
	if err != nil {
		return err
	}

	d.Set("owner", response.Data.Owner)
	d.Set("updated_time", response.Data.LastUpdateTime)
	d.Set("updated_by", response.Data.LastUpdateUser)
	d.Set("result_schema", response.Data.ResultSchema)

	log.Printf("[INFO] Updated Query with guid %s\n", response.Data.QueryID)
	return nil
}

func resourceLaceworkQueryDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting Query with guid %s\n", d.Id())
	_, err := lacework.V2.Query.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Query with guid %s\n", d.Id())
	return nil
}

func importLaceworkQuery(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Query with guid: %s\n", d.Id())

	response, err := lacework.V2.Query.Get(d.Id())
	if err != nil {
		return nil, errors.Errorf(
			"unable to import Lacework resource. Query with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Query found with guid: %s\n", response.Data.QueryID)
	return []*schema.ResourceData{d}, nil
}
