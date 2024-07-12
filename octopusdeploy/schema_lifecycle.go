package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandLifecycle(d *schema.ResourceData) *lifecycles.Lifecycle {
	if d == nil {
		return nil
	}

	name := d.Get("name").(string)

	lifecycle := lifecycles.NewLifecycle(name)
	lifecycle.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		lifecycle.Description = v.(string)
	}

	if v, ok := d.GetOk("phase"); ok {
		lifecycle.Phases = expandPhases(v)
	}

	if v, ok := d.GetOk("release_retention_policy"); ok {
		lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		lifecycle.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("tentacle_retention_policy"); ok {
		lifecycle.TentacleRetentionPolicy = expandRetentionPeriod(v)
	}

	return lifecycle
}

func getLifecycleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "The description of this lifecycle.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"id":   getIDSchema(),
		"name": getNameSchema(true),
		"phase": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getPhaseSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"release_retention_policy": {
			Computed: true,
			DefaultFunc: func() (interface{}, error) {
				return flattenRetentionPeriod(core.NewRetentionPeriod(30, "Days", false)), nil
			},
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_id": getSpaceIDSchema(),
		"tentacle_retention_policy": {
			Computed: true,
			DefaultFunc: func() (interface{}, error) {
				return flattenRetentionPeriod(core.NewRetentionPeriod(30, "Days", false)), nil
			},
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setLifecycle(ctx context.Context, d *schema.ResourceData, lifecycle *lifecycles.Lifecycle) error {
	d.Set("name", lifecycle.Name)
	d.Set("description", lifecycle.Description)
	d.Set("id", lifecycle.GetID())

	if len(lifecycle.SpaceID) > 0 {
		d.Set("space_id", lifecycle.SpaceID)
	}

	if len(lifecycle.Phases) != 0 {
		if err := d.Set("phase", flattenPhases(lifecycle.Phases)); err != nil {
			return fmt.Errorf("error setting phase: %s", err)
		}
	}

	if lifecycle.ReleaseRetentionPolicy != nil {
		if err := d.Set("release_retention_policy", flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy)); err != nil {
			return fmt.Errorf("error setting release_retention_policy: %s", err)
		}
	}

	if lifecycle.TentacleRetentionPolicy != nil {
		if err := d.Set("tentacle_retention_policy", flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy)); err != nil {
			return fmt.Errorf("error setting tentacle_retention_policy: %s", err)
		}
	}

	d.SetId(lifecycle.GetID())
	return nil
}
