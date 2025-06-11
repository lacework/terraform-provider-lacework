package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

func dataSourceLaceworkMetricModule() *schema.Resource {
	return &schema.Resource{
		Read: dataLaceworkMetricModuleRead,

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

func dataLaceworkMetricModuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework      = meta.(*api.Client)
		name          = d.Get("name").(string)
		moduleVersion = d.Get("version").(string)
	)

	metricEvent := api.NewMetricEvent(moduleVersion, name, "lacework-terraform")
	metricEvent.AddFeatureField("lacework_provider_version", version)

	err := lacework.V2.Metrics.Send(metricEvent)
	if err != nil {
		return err
	}

	d.SetId(metricEvent.TraceID)

	return nil
}
