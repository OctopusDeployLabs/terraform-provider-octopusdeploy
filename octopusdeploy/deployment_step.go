package octopusdeploy

import (
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)


func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				"target_roles": &schema.Schema{
					Description: "The roles that this step run against, or runs on behalf of",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"package_requirement": {
					Type:        schema.TypeString,
					Description: "Whether to run this step before or after package acquisition (if possible)",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepPackageRequirement_LetOctopusDecide),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepPackageRequirement_LetOctopusDecide),
						(string)(octopusdeploy.DeploymentStepPackageRequirement_BeforePackageAcquisition),
						(string)(octopusdeploy.DeploymentStepPackageRequirement_AfterPackageAcquisition),
					}),
				},
				"condition": {
					Type:        schema.TypeString,
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepCondition_Success),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepCondition_Success),
						(string)(octopusdeploy.DeploymentStepCondition_Failure),
						(string)(octopusdeploy.DeploymentStepCondition_Always),
						(string)(octopusdeploy.DeploymentStepCondition_Variable),
					}),
				},
				"condition_expression": {
					Type:        schema.TypeString,
					Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
					Optional:    true,
				},
				"start_trigger": {
					Type:        schema.TypeString,
					Description: "Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')",
					Optional:    true,
					Default:     (string)(octopusdeploy.DeploymentStepStartTrigger_StartAfterPrevious),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.DeploymentStepStartTrigger_StartAfterPrevious),
						(string)(octopusdeploy.DeploymentStepStartTrigger_StartWithPrevious),
					}),
				},
				"window_size": {
					Type:        schema.TypeString,
					Description: "The maximum number of targets to deploy to simultaneously",
					Optional:    true,
				},
				"action": getDeploymentActionSchema(),
			},
		},
	}
}

func buildDeploymentStepResource(tfStep map[string]interface{}) octopusdeploy.DeploymentStep {
	step := octopusdeploy.DeploymentStep{
		Name: tfStep["name"].(string),
		PackageRequirement: octopusdeploy.DeploymentStepPackageRequirement(tfStep["package_requirement"].(string)),
		Condition: octopusdeploy.DeploymentStepCondition(tfStep["condition"].(string)),
		StartTrigger: octopusdeploy.DeploymentStepStartTrigger(tfStep["start_trigger"].(string)),
		Properties: map[string]string{},
	}

	targetRoles := tfStep["target_roles"];
	if targetRoles != nil {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	conditionExpression := tfStep["condition_expression"]
	if conditionExpression != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = conditionExpression.(string)
	}

	windowSize := tfStep["window_size"]
	if windowSize != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = windowSize.(string)
	}

	if attr, ok := tfStep["action"]; ok {
		for _, tfAction := range attr.(*schema.Set).List() {
			action := buildDeploymentActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step;
}

