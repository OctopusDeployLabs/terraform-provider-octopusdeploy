package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				"target_roles": {
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
				"action":                          getDeploymentActionSchema(),
				"manual_intervention_action":      getManualInterventionActionSchema(),
				"apply_terraform_action":          getApplyTerraformActionSchema(),
				"deploy_package_action":           getDeployPackageAction(),
				"deploy_windows_service_action":   getDeployWindowsServiceActionSchema(),
				"run_script_action":               getRunScriptActionSchema(),
				"run_kubectl_script_action":       getRunRunKubectlScriptSchema(),
				"deploy_kubernetes_secret_action": getDeployKubernetesSecretActionSchema(),
			},
		},
	}
}

func buildDeploymentStepResource(tfStep map[string]interface{}) octopusdeploy.DeploymentStep {
	step := octopusdeploy.DeploymentStep{
		Name:               tfStep["name"].(string),
		PackageRequirement: octopusdeploy.DeploymentStepPackageRequirement(tfStep["package_requirement"].(string)),
		Condition:          octopusdeploy.DeploymentStepCondition(tfStep["condition"].(string)),
		StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(tfStep["start_trigger"].(string)),
		Properties:         map[string]string{},
	}

	targetRoles := tfStep["target_roles"]
	if targetRoles != nil {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	conditionExpression := tfStep["condition_expression"]
	if conditionExpression != nil {
		step.Properties["Octopus.Action.ConditionVariableExpression"] = conditionExpression.(string)
	}

	windowSize := tfStep["window_size"]
	if windowSize != nil {
		step.Properties["Octopus.Action.MaxParallelism"] = windowSize.(string)
	}

	if attr, ok := tfStep["action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeploymentActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["manual_intervention_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildManualInterventionActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["apply_terraform_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildApplyTerraformActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["deploy_package_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployPackageActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["deploy_windows_service_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployWindowsServiceActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["run_script_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildRunScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["run_kubectl_script_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildRunKubectlScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if attr, ok := tfStep["deploy_kubernetes_secret_action"]; ok {
		for _, tfAction := range attr.([]interface{}) {
			action := buildDeployKubernetesSecretActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step
}
