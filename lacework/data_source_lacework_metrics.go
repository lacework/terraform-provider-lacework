package lacework

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkMetrics() *schema.Resource {
	return &schema.Resource{
		Read: dataLaceworkMetricsRead,

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
				Default:     "lacework-terraform",
				Description: "The name of the dataset",
			},
			"trace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataLaceworkMetricsRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework       = meta.(*api.Client)
		dataset        = d.Get("dataset").(string)
		name           = d.Get("name").(string)
		module_version = d.Get("version").(string)
	)

	honeycombEvent := api.NewHoneyvent(
		module_version,
		name,
		dataset)
	resp, err := lacework.V2.Metrics.Send(honeycombEvent)
	if err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	d.Set("trace_id", resp.Data[0].TraceID)

	return nil
}
