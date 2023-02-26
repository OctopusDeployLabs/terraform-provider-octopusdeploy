package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPhase(tfPhase map[string]interface{}) *lifecycles.Phase {
	if tfPhase == nil {
		return nil
	}

	name := tfPhase["name"].(string)

	phase := lifecycles.NewPhase(name)

	if v, ok := tfPhase["automatic_deployment_targets"]; ok {
		phase.AutomaticDeploymentTargets = getSliceFromTerraformTypeList(v)
	}

	if v, ok := tfPhase["is_optional_phase"]; ok {
		phase.IsOptionalPhase = v.(bool)
	}

	if v, ok := tfPhase["minimum_environments_before_promotion"]; ok {
		if n, isInt32 := v.(int32); isInt32 {
			phase.MinimumEnvironmentsBeforePromotion = n
		} else {
			phase.MinimumEnvironmentsBeforePromotion = int32(v.(int))
		}
	}

	if v, ok := tfPhase["optional_deployment_targets"]; ok {
		phase.OptionalDeploymentTargets = getSliceFromTerraformTypeList(v)
	}

	if phase.AutomaticDeploymentTargets == nil {
		phase.AutomaticDeploymentTargets = []string{}
	}

	if phase.OptionalDeploymentTargets == nil {
		phase.OptionalDeploymentTargets = []string{}
	}

	if v, ok := tfPhase["release_retention_policy"]; ok {
		retentionPolicy := v.([]interface{})
		if len(retentionPolicy) == 1 {
			phase.ReleaseRetentionPolicy = expandRetentionPeriod(retentionPolicy[0].(map[string]interface{}))
		}
	}

	if v, ok := tfPhase["tentacle_retention_policy"]; ok {
		retentionPolicy := v.([]interface{})
		if len(retentionPolicy) == 1 {
			phase.TentacleRetentionPolicy = expandRetentionPeriod(retentionPolicy[0].(map[string]interface{}))
		}
	}

	return phase
}

func flattenPhases(phases []*lifecycles.Phase) []interface{} {
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
			p["release_retention_policy"] = flattenRetentionPeriod(phase.ReleaseRetentionPolicy)
		}
		if phase.TentacleRetentionPolicy != nil {
			p["tentacle_retention_policy"] = flattenRetentionPeriod(phase.TentacleRetentionPolicy)
		}
		flattenedPhases = append(flattenedPhases, p)
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
