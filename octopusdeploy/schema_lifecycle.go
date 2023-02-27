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

	if v, ok := d.GetOk("release_retention_policy"); ok {
		lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		lifecycle.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("tentacle_retention_policy"); ok {
		lifecycle.TentacleRetentionPolicy = expandRetentionPeriod(v)
	}

	if v, ok := d.GetOk("release_retention_policy"); ok {
		retentionPeriod := v.([]interface{})
		if len(retentionPeriod) == 1 {
			lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(retentionPeriod[0].(map[string]interface{}))
		}
	}

	if attr, ok := d.GetOk("phase"); ok {
		tfPhases := attr.([]interface{})

		for _, tfPhase := range tfPhases {
			phase := expandPhase(tfPhase.(map[string]interface{}))
			lifecycle.Phases = append(lifecycle.Phases, phase)
		}
	}

	return lifecycle
}

func flattenLifecycle(lifecycle *lifecycles.Lifecycle) map[string]interface{} {
	if lifecycle == nil {
		return nil
	}

	return map[string]interface{}{
		"description":               lifecycle.Description,
		"id":                        lifecycle.GetID(),
		"name":                      lifecycle.Name,
		"phase":                     flattenPhases(lifecycle.Phases),
		"space_id":                  lifecycle.SpaceID,
		"release_retention_policy":  flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy),
		"tentacle_retention_policy": flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy),
	}
}

func getLifecycleDataSchema() map[string]*schema.Schema {
	dataSchema := getLifecycleSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids": getQueryIDs(),
		"lifecycles": {
			Computed:    true,
			Description: "A list of lifecycles that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
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
