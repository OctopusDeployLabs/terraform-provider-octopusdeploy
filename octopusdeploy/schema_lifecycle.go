package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func expandRetentionPeriod(d *schema.ResourceData, key string) *octopusdeploy.RetentionPeriod {
	v, ok := d.GetOk(key)
	if ok {
		retentionPeriod := v.([]interface{})
		if len(retentionPeriod) == 1 {
			tfRetentionItem := retentionPeriod[0].(map[string]interface{})
			retention := octopusdeploy.RetentionPeriod{
				QuantityToKeep:    int32(tfRetentionItem["quantity_to_keep"].(int)),
				ShouldKeepForever: tfRetentionItem["should_keep_forever"].(bool),
				Unit:              tfRetentionItem["unit"].(string),
			}
			return &retention
		}
	}

	return nil
}

func expandPhase(tfPhase map[string]interface{}) octopusdeploy.Phase {
	phase := octopusdeploy.Phase{
		AutomaticDeploymentTargets:         getSliceFromTerraformTypeList(tfPhase["automatic_deployment_targets"]),
		IsOptionalPhase:                    tfPhase["is_optional_phase"].(bool),
		MinimumEnvironmentsBeforePromotion: int32(tfPhase["minimum_environments_before_promotion"].(int)),
		Name:                               tfPhase["name"].(string),
		OptionalDeploymentTargets:          getSliceFromTerraformTypeList(tfPhase["optional_deployment_targets"]),
	}

	if phase.AutomaticDeploymentTargets == nil {
		phase.AutomaticDeploymentTargets = []string{}
	}
	if phase.OptionalDeploymentTargets == nil {
		phase.OptionalDeploymentTargets = []string{}
	}

	return phase
}

func flattenLifecycle(lifecycle *octopusdeploy.Lifecycle) map[string]interface{} {
	if lifecycle == nil {
		return nil
	}

	return map[string]interface{}{
		"description":               lifecycle.Description,
		"id":                        lifecycle.GetID(),
		"name":                      lifecycle.Name,
		"phases":                    lifecycle.Phases,
		"space_id":                  lifecycle.SpaceID,
		"release_retention_policy":  flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy),
		"tentacle_retention_policy": flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy),
	}
}

func flattenPhases(phases []octopusdeploy.Phase) []interface{} {
	flattenedPhases := make([]interface{}, 0)
	for _, phase := range phases {
		p := make(map[string]interface{})
		p["automatic_deployment_targets"] = flattenArray(phase.AutomaticDeploymentTargets)
		p["id"] = phase.ID
		p["is_optional_phase"] = phase.IsOptionalPhase
		p["minimum_environments_before_promotion"] = int(phase.MinimumEnvironmentsBeforePromotion)
		p["name"] = phase.Name
		p["optional_deployment_targets"] = flattenArray(phase.OptionalDeploymentTargets)
		if phase.ReleaseRetentionPolicy != nil {
			p["release_retention_policy"] = flattenRetentionPeriod(*phase.ReleaseRetentionPolicy)
		}
		if phase.TentacleRetentionPolicy != nil {
			p["tentacle_retention_policy"] = flattenRetentionPeriod(*phase.TentacleRetentionPolicy)
		}
		flattenedPhases = append(flattenedPhases, p)
	}
	return flattenedPhases
}

func flattenRetentionPeriod(r octopusdeploy.RetentionPeriod) []interface{} {
	retentionPeriod := make(map[string]interface{})
	retentionPeriod["quantity_to_keep"] = int(r.QuantityToKeep)
	retentionPeriod["should_keep_forever"] = r.ShouldKeepForever
	retentionPeriod["unit"] = r.Unit
	return []interface{}{retentionPeriod}
}

func getLifecycleDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"phases": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getPhaseDataSchema()},
			Type:     schema.TypeList,
		},
		"release_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Type:     schema.TypeList,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tentacle_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Type:     schema.TypeList,
		},
	}
}

func getLifecycleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"phases": {
			Elem:     &schema.Resource{Schema: getPhaseSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"release_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tentacle_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func getPhaseDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"automatic_deployment_targets": {
			Computed:    true,
			Description: "Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Type:        schema.TypeList,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_optional_phase": {
			Computed:    true,
			Description: "If false a release must be deployed to this phase before it can be deployed to the next phase.",
			Type:        schema.TypeBool,
		},
		"minimum_environments_before_promotion": {
			Computed:    true,
			Description: "The number of units required before a release can enter the next phase. If 0, all environments are required.",
			Type:        schema.TypeInt,
		},
		"name": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"optional_deployment_targets": {
			Computed:    true,
			Description: "Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"release_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodDataSchema()},
			Type:     schema.TypeList,
		},
		"tentacle_retention_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getRetentionPeriodDataSchema()},
			Type:     schema.TypeList,
		},
	}
}

func getPhaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"automatic_deployment_targets": {
			Description: "Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_optional_phase": {
			Default:     false,
			Description: "If false a release must be deployed to this phase before it can be deployed to the next phase.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"minimum_environments_before_promotion": {
			Default:     0,
			Description: "The number of units required before a release can enter the next phase. If 0, all environments are required.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"optional_deployment_targets": {
			Description: "Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"release_retention_policy": {
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tentacle_retention_policy": {
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func getRetentionPeriodDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"quantity_to_keep": {
			Computed:    true,
			Description: "The number of days/releases to keep. If 0 all are kept.",
			Type:        schema.TypeInt,
		},
		"should_keep_forever": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"unit": {
			Computed:    true,
			Description: "The unit of quantity_to_keep.",
			Type:        schema.TypeString,
		},
	}
}

func getRetentionPeriodSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"quantity_to_keep": {
			Default:     30,
			Description: "The number of days/releases to keep. If 0 all are kept.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"should_keep_forever": {
			Default:  false,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"unit": {
			Default:     octopusdeploy.RetentionUnitDays,
			Description: "The unit of quantity_to_keep.",
			Optional:    true,
			Type:        schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				octopusdeploy.RetentionUnitDays,
				octopusdeploy.RetentionUnitItems,
			}, false)),
		},
	}
}

func setLifecycle(ctx context.Context, d *schema.ResourceData, lifecycle *octopusdeploy.Lifecycle) {
	d.Set("description", lifecycle.Description)
	d.Set("id", lifecycle.GetID())
	d.Set("name", lifecycle.Name)
	d.Set("phases", flattenPhases(lifecycle.Phases))
	d.Set("space_id", lifecycle.SpaceID)
	d.Set("release_retention_policy", flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy))
	d.Set("tentacle_retention_policy", flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy))

	d.SetId(lifecycle.GetID())
}
