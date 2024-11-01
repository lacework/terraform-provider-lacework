package lacework

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
	"github.com/pkg/errors"
)

func resourceLaceworkQuery() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkQueryCreate,
		Read:   resourceLaceworkQueryRead,
		Update: resourceLaceworkQueryUpdate,
		Delete: resourceLaceworkQueryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: importLaceworkQuery,
		},

		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the query",
			},
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

func resourceLaceworkQueryCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
	)

	query := api.NewQuery{
		QueryID:   d.Get("query_id").(string),
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
		return resourceNotFound(d, err)
	}

	d.Set("query", response.Data.QueryText)
	d.Set("query_id", response.Data.QueryID)
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

	if d.HasChange("query_id") {
		return errors.New("unable to change ID of an existing query")
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

func importLaceworkQuery(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Query with guid: %s\n", d.Id())

	response, err := lacework.V2.Query.Get(d.Id())
	if err != nil {
		return nil, fmt.Errorf(
			"unable to import Lacework resource. Query with guid '%s' was not found",
			d.Id(),
		)
	}
	log.Printf("[INFO] Query found with guid: %s\n", response.Data.QueryID)
	return []*schema.ResourceData{d}, nil
}
