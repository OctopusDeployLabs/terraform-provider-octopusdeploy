package octopusdeploy

import (
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentStep(flattenedStep map[string]interface{}) *octopusdeploy.DeploymentStep {
	name := flattenedStep["name"].(string)
	step := octopusdeploy.NewDeploymentStep(name)

	// properties MUST be serialized first
	if properties, ok := flattenedStep["properties"]; ok {
		step.Properties = expandProperties(properties)
	}

	if condition, ok := flattenedStep["condition"]; ok {
		step.Condition = octopusdeploy.DeploymentStepConditionType(condition.(string))
	}

	if conditionExpression, ok := flattenedStep["condition_expression"]; ok {
		step.Properties["Octopus.Step.ConditionVariableExpression"] = octopusdeploy.NewPropertyValue(conditionExpression.(string), false)
	}

	if packageRequirement, ok := flattenedStep["package_requirement"]; ok {
		step.PackageRequirement = octopusdeploy.DeploymentStepPackageRequirement(packageRequirement.(string))
	}

	if startTrigger, ok := flattenedStep["start_trigger"]; ok {
		step.StartTrigger = octopusdeploy.DeploymentStepStartTrigger(startTrigger.(string))
	}

	if targetRoles, ok := flattenedStep["target_roles"]; ok {
		step.Properties["Octopus.Action.TargetRoles"] = octopusdeploy.NewPropertyValue(strings.Join(getSliceFromTerraformTypeList(targetRoles), ","), false)
	}

	if windowSize, ok := flattenedStep["window_size"]; ok {
		step.Properties["Octopus.Action.MaxParallelism"] = octopusdeploy.NewPropertyValue(windowSize.(string), false)
	}

	if v, ok := flattenedStep["action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["manual_intervention_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandManualInterventionAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["apply_terraform_template_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandApplyTerraformTemplateAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["deploy_package_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployPackageAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["deploy_windows_service_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployWindowsServiceAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["run_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandRunScriptAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["run_kubectl_script_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandRunKubectlScriptAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
		}
	}

	if v, ok := flattenedStep["deploy_kubernetes_secret_action"]; ok {
		for _, tfAction := range v.([]interface{}) {
			action := expandDeployKubernetesSecretAction(tfAction.(map[string]interface{}))
			step.Actions = append(step.Actions, *action)
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
		flattenedDeploymentSteps[key]["properties"] = flattenProperties(deploymentStep.Properties)
		flattenedDeploymentSteps[key]["start_trigger"] = deploymentStep.StartTrigger

		for propertyName, propertyValue := range deploymentStep.Properties {
			switch propertyName {
			case "Octopus.Action.TargetRoles":
				flattenedDeploymentSteps[key]["target_roles"] = strings.Split(propertyValue.Value, ",")
			case "Octopus.Action.MaxParallelism":
				flattenedDeploymentSteps[key]["window_size"] = propertyValue.Value
			case "Octopus.Step.ConditionVariableExpression":
				flattenedDeploymentSteps[key]["condition_expression"] = propertyValue.Value
			}
		}

		for i := range deploymentStep.Actions {
			switch deploymentStep.Actions[i].ActionType {
			case "Octopus.KubernetesDeploySecret":
				flattenedDeploymentSteps[key]["deploy_kubernetes_secret_action"] = []interface{}{flattenDeployKubernetesSecretAction(&deploymentStep.Actions[i])}
			case "Octopus.KubernetesRunScript":
				flattenedDeploymentSteps[key]["run_kubectl_script_action"] = []interface{}{flattenKubernetesRunScriptAction(&deploymentStep.Actions[i])}
			case "Octopus.Manual":
				flattenedDeploymentSteps[key]["manual_intervention_action"] = []interface{}{flattenManualInterventionAction(&deploymentStep.Actions[i])}
			case "Octopus.Script":
				flattenedDeploymentSteps[key]["run_script_action"] = []interface{}{flattenRunScriptAction(&deploymentStep.Actions[i])}
			case "Octopus.TentaclePackage":
				flattenedDeploymentSteps[key]["deploy_package_action"] = []interface{}{flattenDeployPackageAction(&deploymentStep.Actions[i])}
			case "Octopus.TerraformApply":
				flattenedDeploymentSteps[key]["apply_terraform_template_action"] = []interface{}{flattenApplyTerraformTemplateAction(&deploymentStep.Actions[i])}
			case "Octopus.WindowsService":
				flattenedDeploymentSteps[key]["deploy_windows_service_action"] = []interface{}{flattenDeployWindowsServiceAction(&deploymentStep.Actions[i])}
			default:
				flattenedDeploymentSteps[key]["action"] = []interface{}{flattenDeploymentAction(&deploymentStep.Actions[i])}
			}
		}
	}

	return flattenedDeploymentSteps
}

func getDeploymentStepSchema() *schema.Schema {
	return &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action":                          getDeploymentActionSchema(),
				"apply_terraform_template_action": getApplyTerraformTemplateActionSchema(),
				"condition": {
					Default:     "Success",
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
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
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"AfterPackageAcquisition",
						"BeforePackageAcquisition",
						"LetOctopusDecide",
					}, false)),
				},
				"properties": {
					Computed: true,
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
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"StartAfterPrevious",
						"StartWithPrevious",
					}, false)),
				},
				"target_roles": {
					Computed:    true,
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
