package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func getPhaseDataSchema() map[string]*schema.Schema {
	dataSchema := getPhaseSchema()
	setDataSchema(&dataSchema)

	return dataSchema
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
