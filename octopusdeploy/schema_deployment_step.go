package octopusdeploy

import (
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentStep(tfStep map[string]interface{}) *octopusdeploy.DeploymentStep {
	name := tfStep["name"].(string)
	step := octopusdeploy.NewDeploymentStep(name)

	// properties MUST be serialized first
	if properties, ok := tfStep["properties"]; ok {
		step.Properties = expandProperties(properties)
	}

	if condition, ok := tfStep["condition"]; ok {
		step.Condition = octopusdeploy.DeploymentStepConditionType(condition.(string))
	}

	if conditionExpression, ok := tfStep["condition_expression"]; ok {
		step.Properties["Octopus.Step.ConditionVariableExpression"] = conditionExpression.(string)
	}

	if packageRequirement, ok := tfStep["package_requirement"]; ok {
		step.PackageRequirement = octopusdeploy.DeploymentStepPackageRequirement(packageRequirement.(string))
	}

	if startTrigger, ok := tfStep["start_trigger"]; ok {
		step.StartTrigger = octopusdeploy.DeploymentStepStartTrigger(startTrigger.(string))
	}

	if targetRoles, ok := tfStep["target_roles"]; ok {
		step.Properties["Octopus.Action.TargetRoles"] = strings.Join(getSliceFromTerraformTypeList(targetRoles), ",")
	}

	if windowSize, ok := tfStep["window_size"]; ok {
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
			action := expandManualInterventionAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["apply_terraform_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandApplyTerraformAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_package_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployPackageAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_windows_service_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployWindowsServiceAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["run_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandRunScriptAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["run_kubectl_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandRunKubectlScriptAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	if v, ok := tfStep["deploy_kubernetes_secret_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployKubernetesSecretAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, action)
		}
	}

	return step
}

func flattenDeploymentSteps(deploymentSteps []octopusdeploy.DeploymentStep) []map[string]interface{} {
	if deploymentSteps == nil {
		return nil
	}

	var flattenedDeploymentSteps = make([]map[string]interface{}, len(deploymentSteps))
	for key, deploymentStep := range deploymentSteps {
		flattenedDeploymentSteps[key] = map[string]interface{}{}
		flattenedDeploymentSteps[key]["condition"] = deploymentStep.Condition
		flattenedDeploymentSteps[key]["id"] = deploymentStep.ID
		flattenedDeploymentSteps[key]["name"] = deploymentStep.Name
		flattenedDeploymentSteps[key]["package_requirement"] = deploymentStep.PackageRequirement
		flattenedDeploymentSteps[key]["properties"] = deploymentStep.Properties
		flattenedDeploymentSteps[key]["start_trigger"] = deploymentStep.StartTrigger

		for propertyName, propertyValue := range deploymentStep.Properties {
			switch propertyName {
			case "Octopus.Action.TargetRoles":
				flattenedDeploymentSteps[key]["target_roles"] = strings.Split(propertyValue, ",")
			case "Octopus.Action.MaxParallelism":
				flattenedDeploymentSteps[key]["window_size"] = propertyValue
			case "Octopus.Step.ConditionVariableExpression":
				flattenedDeploymentSteps[key]["condition_expression"] = propertyValue
			}
		}

		for _, action := range deploymentStep.Actions {
			switch action.ActionType {
			case "Octopus.KubernetesDeploySecret":
				flattenedDeploymentSteps[key]["deploy_kubernetes_secret_action"] = []interface{}{flattenDeployKubernetesSecretAction(action)}
			case "Octopus.KubernetesRunScript":
				flattenedDeploymentSteps[key]["run_kubectl_script_action"] = []interface{}{flattenKubernetesRunScriptAction(action)}
			case "Octopus.Manual":
				flattenedDeploymentSteps[key]["manual_intervention_action"] = []interface{}{flattenManualInterventionAction(action)}
			case "Octopus.Script":
				flattenedDeploymentSteps[key]["run_script_action"] = []interface{}{flattenRunScriptAction(action)}
			case "Octopus.TentaclePackage":
				flattenedDeploymentSteps[key]["deploy_package_action"] = []interface{}{flattenDeployPackageAction(action)}
			case "Octopus.TerraformApply":
				flattenedDeploymentSteps[key]["apply_terraform_action"] = []interface{}{flattenApplyTerraformAction(action)}
			case "Octopus.WindowsService":
				flattenedDeploymentSteps[key]["deploy_windows_service_action"] = []interface{}{flattenDeployWindowsServiceAction(action)}
			default:
				flattenedDeploymentSteps[key]["action"] = []interface{}{flattenDeploymentAction(action)}
			}
		}
	}

	return flattenedDeploymentSteps
}

func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
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
					Computed:    true,
					Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"deploy_kubernetes_secret_action": getDeployKubernetesSecretActionSchema(),
				"deploy_package_action":           getDeployPackageActionSchema(),
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
				"properties": {
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Type:     schema.TypeMap,
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
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Type:        schema.TypeList,
				},
				"window_size": {
					Description: "The maximum number of targets to deploy to simultaneously",
					Optional:    true,
					Type:        schema.TypeString,
				},
			},
		},
		Optional: true,
		Type:     schema.TypeList,
	}
}
