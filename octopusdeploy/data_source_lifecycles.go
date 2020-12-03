package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLifecycles() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing lifecycles.",
		ReadContext: dataSourceLifecyclesRead,
		Schema:      getLifecycleDataSchema(),
	}
}

func dataSourceLifecyclesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.LifecyclesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	lifecycles, err := client.Lifecycles.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedLifecycles := []interface{}{}
	for _, lifecycle := range lifecycles.Items {
		flattenedLifecycles = append(flattenedLifecycles, flattenLifecycle(lifecycle))
	}

	d.Set("lifecycles", flattenedLifecycles)
	d.SetId("Lifecycles " + time.Now().UTC().String())

	return nil
}
