package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkMetrics() *schema.Resource {
	return &schema.Resource{
		Create: dataLaceworkMetricsCreate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the module",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the module",
			},
			"dataset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the dataset",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func dataLaceworkMetricsCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		dataset  = d.Get("dataset").(string)
		name     = d.Get("name").(string)
		version  = d.Get("version").(string)
	)

	if dataset == "" {
		dataset = "lacework-terraform"
	}

	honeycombEvent := api.Honeyvent{
		Version: version,
		Feature: name,
		Dataset: dataset,
	}

	_, err := lacework.V2.Metrics.Send(honeycombEvent)
	if err != nil {
		return err
	}

	return nil
}
