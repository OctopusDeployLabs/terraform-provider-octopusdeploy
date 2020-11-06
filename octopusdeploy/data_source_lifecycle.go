package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLifecycle() *schema.Resource {
	dataSourceLifecycleSchema := map[string]*schema.Schema{
		constDescription: &schema.Schema{
			Optional: true,
			Type:     schema.TypeString,
		},
		constName: &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
		constPhase: {
			Elem:     phaseSchema(),
			Optional: true,
			Set:      schema.HashResource(phaseSchema()),
			Type:     schema.TypeSet,
		},
		constReleaseRetentionPolicy: getRetentionPeriodSchema(),
		constSpaceID: {
			Type:     schema.TypeString,
			Computed: true,
		},
		constTentacleRetentionPolicy: getRetentionPeriodSchema(),
	}

	return &schema.Resource{
		ReadContext: dataSourceLifecycleRead,
		Schema:      dataSourceLifecycleSchema,
	}
}

func dataSourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get(constName).(string)

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
			logResource(constLifecycle, m)

			flattenLifecycle(ctx, d, lifecycle)

			return nil
		}
	}

	return nil
}
