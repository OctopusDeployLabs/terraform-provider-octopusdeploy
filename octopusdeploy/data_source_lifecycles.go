package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLifecycle() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLifecycleRead,
		Schema:      getLifecycleDataSchema(),
	}
}

func dataSourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	lifecycles, err := client.Lifecycles.GetByPartialName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(lifecycles) == 0 {
		return nil
	}

	// NOTE: two or more lifecycles can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, lifecycle := range lifecycles {
		if lifecycle.Name == name {
			flattenLifecycle(ctx, d, lifecycle)
			return nil
		}
	}

	return nil
}
