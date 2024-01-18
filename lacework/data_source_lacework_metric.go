package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkMetric() *schema.Resource {
	return &schema.Resource{
		Read: dataLaceworkMetricRead,

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
		},
	}
}

func dataLaceworkMetricRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		name          = d.Get("name").(string)
		moduleVersion = d.Get("version").(string)
	)

	honeycombEvent := api.NewHoneyvent(moduleVersion, name, "lacework-terraform")
	resp, err := lacework.V2.Metrics.Send(honeycombEvent)
	if err != nil {
		return err
	}

	d.SetId(resp.Data[0].TraceID)

	return nil
}
