package octopusdeploy

import (
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentStep(tfStep map[string]interface{}) octopusdeploy.DeploymentStep {
	step := octopusdeploy.DeploymentStep{
		Name:               tfStep["name"].(string),
		PackageRequirement: octopusdeploy.DeploymentStepPackageRequirement(tfStep["package_requirement"].(string)),
		Condition:          octopusdeploy.DeploymentStepConditionType(tfStep["condition"].(string)),
		StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(tfStep["start_trigger"].(string)),
		Properties:         map[string]string{},
	}

	targetRoles := tfStep["target_roles"]
	if targetRoles != nil {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	conditionExpression := tfStep["condition_expression"]
	if conditionExpression != nil {
		step.Properties["Octopus.Step.ConditionVariableExpression"] = conditionExpression.(string)
	}

	windowSize := tfStep["window_size"]
	if windowSize != nil {
		step.Properties["Octopus.Action.MaxParallelism"] = windowSize.(string)
	}

	if v, ok := tfStep["action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeploymentAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["manual_intervention_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildManualInterventionActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["apply_terraform_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildApplyTerraformActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_package_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildDeployPackageActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_windows_service_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildDeployWindowsServiceActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["run_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildRunScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["run_kubectl_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildRunKubectlScriptActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_kubernetes_secret_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := buildDeployKubernetesSecretActionResource(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step
}

func flattenDeploymentSteps(deploymentSteps []octopusdeploy.DeploymentStep) []map[string]interface{} {
	var flattenedDeploymentSteps = make([]map[string]interface{}, len(deploymentSteps))
	for key, deploymentStep := range deploymentSteps {
		flattenedDeploymentSteps[key] = map[string]interface{}{
			"actions":             flattenDeploymentActions(deploymentStep.Actions),
			"condition":           deploymentStep.Condition,
			"id":                  deploymentStep.ID,
			"name":                deploymentStep.Name,
			"package_requirement": deploymentStep.PackageRequirement,
			"properties":          deploymentStep.Properties,
			"start_trigger":       deploymentStep.StartTrigger,
		}
	}

	return flattenedDeploymentSteps
}

func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action":                 getDeploymentActionSchema(),
				"apply_terraform_action": getApplyTerraformActionSchema(),
				"condition": {
					Default:     "Success",
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						"Always",
						"Failure",
						"Success",
						"Variable",
					}, false)),
				},
				"condition_expression": {
					Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"deploy_kubernetes_secret_action": getDeployKubernetesSecretActionSchema(),
				"deploy_package_action":           getDeployPackageAction(),
				"deploy_windows_service_action":   getDeployWindowsServiceActionSchema(),
				"id":                              getIDSchema(),
				"manual_intervention_action":      getManualInterventionActionSchema(),
				"name":                            getNameSchema(true),
				"package_requirement": {
					Default:     "LetOctopusDecide",
					Description: "Whether to run this step before or after package acquisition (if possible)",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						"AfterPackageAcquisition",
						"BeforePackageAcquisition",
						"LetOctopusDecide",
					}, false)),
				},
				"run_kubectl_script_action": getRunKubectlScriptSchema(),
				"run_script_action":         getRunScriptActionSchema(),
				"start_trigger": {
					Default:     "StartAfterPrevious",
					Description: "Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						"StartAfterPrevious",
						"StartWithPrevious",
					}, false)),
				},
				"target_roles": {
					Description: "The roles that this step run against, or runs on behalf of",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
					Type:     schema.TypeList,
				},
				"window_size": {
					Description: "The maximum number of targets to deploy to simultaneously",
					Optional:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}
}
