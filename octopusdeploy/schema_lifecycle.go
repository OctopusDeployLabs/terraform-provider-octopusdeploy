package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandLifecycle(d *schema.ResourceData) *octopusdeploy.Lifecycle {
	name := d.Get("name").(string)

	lifecycle := octopusdeploy.NewLifecycle(name)
	lifecycle.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		lifecycle.Description = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		lifecycle.SpaceID = v.(string)
	}

	releaseRetentionPolicy := expandRetentionPeriod(d, "release_retention_policy")
	if releaseRetentionPolicy != nil {
		lifecycle.ReleaseRetentionPolicy = *releaseRetentionPolicy
	}

	tentacleRetentionPolicy := expandRetentionPeriod(d, "tentacle_retention_policy")
	if tentacleRetentionPolicy != nil {
		lifecycle.TentacleRetentionPolicy = *tentacleRetentionPolicy
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

func flattenLifecycle(lifecycle *octopusdeploy.Lifecycle) map[string]interface{} {
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
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"phase": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getPhaseSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"release_retention_policy": {
			Computed: true,
			DefaultFunc: func() (interface{}, error) {
				return flattenRetentionPeriod(octopusdeploy.RetentionPeriod{
					Unit:           "Days",
					QuantityToKeep: 30,
				}), nil
			},
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			MaxItems: 1,
			MinItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_id": getSpaceIDSchema(),
		"tentacle_retention_policy": {
			Computed: true,
			DefaultFunc: func() (interface{}, error) {
				return flattenRetentionPeriod(octopusdeploy.RetentionPeriod{
					Unit:           "Days",
					QuantityToKeep: 30,
				}), nil
			},
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			MaxItems: 1,
			MinItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setLifecycle(ctx context.Context, d *schema.ResourceData, lifecycle *octopusdeploy.Lifecycle) error {
	d.Set("description", lifecycle.Description)
	d.Set("id", lifecycle.GetID())
	d.Set("name", lifecycle.Name)
	d.Set("space_id", lifecycle.SpaceID)

	if err := d.Set("phase", flattenPhases(lifecycle.Phases)); err != nil {
		return fmt.Errorf("error setting phase: %s", err)
	}

	if err := d.Set("release_retention_policy", flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy)); err != nil {
		return fmt.Errorf("error setting release_retention_policy: %s", err)
	}

	if err := d.Set("tentacle_retention_policy", flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy)); err != nil {
		return fmt.Errorf("error setting tentacle_retention_policy: %s", err)
	}

	d.SetId(lifecycle.GetID())

	return nil
}
