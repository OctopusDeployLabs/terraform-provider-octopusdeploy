package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPhase(flattenedPhase interface{}) *lifecycles.Phase {
	if flattenedPhase == nil {
		return nil
	}

	if _, ok := flattenedPhase.(map[string]interface{}); !ok {
		return nil
	}

	flattenedValues := flattenedPhase.(map[string]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	name := flattenedValues["name"].(string)
	phase := lifecycles.NewPhase(name)

	if v, ok := flattenedValues["automatic_deployment_targets"]; ok {
		phase.AutomaticDeploymentTargets = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedValues["is_optional_phase"]; ok {
		phase.IsOptionalPhase = v.(bool)
	}

	if v, ok := flattenedValues["minimum_environments_before_promotion"]; ok {
		if n, isInt32 := v.(int32); isInt32 {
			phase.MinimumEnvironmentsBeforePromotion = n
		} else {
			phase.MinimumEnvironmentsBeforePromotion = int32(v.(int))
		}
	}

	if v, ok := flattenedValues["optional_deployment_targets"]; ok {
		phase.OptionalDeploymentTargets = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedValues["release_retention_policy"]; ok {
		phase.ReleaseRetentionPolicy = expandRetentionPeriod(v)
	}

	if v, ok := flattenedValues["tentacle_retention_policy"]; ok {
		phase.TentacleRetentionPolicy = expandRetentionPeriod(v)
	}

	if phase.AutomaticDeploymentTargets == nil {
		phase.AutomaticDeploymentTargets = []string{}
	}

	if phase.OptionalDeploymentTargets == nil {
		phase.OptionalDeploymentTargets = []string{}
	}

	return phase
}

func expandPhases(flattenedPhases interface{}) []*lifecycles.Phase {
	if flattenedPhases == nil {
		return nil
	}

	if _, ok := flattenedPhases.([]interface{}); !ok {
		return nil
	}

	flattenedValues := flattenedPhases.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	phases := []*lifecycles.Phase{}
	for _, flattenedValues := range flattenedValues {
		phases = append(phases, expandPhase(flattenedValues))
	}
	return phases
}

func flattenPhase(phase *lifecycles.Phase) interface{} {
	if phase == nil {
		return nil
	}

	flattenedPhase := make(map[string]interface{})
	flattenedPhase["automatic_deployment_targets"] = flattenArray(phase.AutomaticDeploymentTargets)
	flattenedPhase["id"] = phase.ID
	flattenedPhase["is_optional_phase"] = phase.IsOptionalPhase
	flattenedPhase["minimum_environments_before_promotion"] = int(phase.MinimumEnvironmentsBeforePromotion)
	flattenedPhase["name"] = phase.Name
	flattenedPhase["optional_deployment_targets"] = flattenArray(phase.OptionalDeploymentTargets)
	if phase.ReleaseRetentionPolicy != nil {
		flattenedPhase["release_retention_policy"] = flattenRetentionPeriod(phase.ReleaseRetentionPolicy)
	}
	if phase.TentacleRetentionPolicy != nil {
		flattenedPhase["tentacle_retention_policy"] = flattenRetentionPeriod(phase.TentacleRetentionPolicy)
	}

	return flattenedPhase
}

func flattenPhases(phases []*lifecycles.Phase) []interface{} {
	flattenedPhases := make([]interface{}, 0)
	for _, phase := range phases {
		flattenedPhases = append(flattenedPhases, flattenPhase(phase))
	}
	return flattenedPhases
}

func getPhaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"automatic_deployment_targets": {
			Description: "Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"id": getIDSchema(),
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
		"name": getNameSchema(true),
		"optional_deployment_targets": {
			Description: "Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"release_retention_policy": {
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			MaxItems: 1,
			Type:     schema.TypeList,
		},
		"tentacle_retention_policy": {
			Elem:     &schema.Resource{Schema: getRetentionPeriodSchema()},
			Optional: true,
			MaxItems: 1,
			Type:     schema.TypeList,
		},
	}
}
