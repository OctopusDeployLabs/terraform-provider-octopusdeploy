package octopusdeploy

import (
	"log"
	"sort"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentStep(flattenedStep map[string]interface{}) *deployments.DeploymentStep {
	name := flattenedStep["name"].(string)
	step := deployments.NewDeploymentStep(name)

	// properties MUST be serialized first
	if properties, ok := flattenedStep["properties"]; ok {
		step.Properties = expandProperties(properties)
	}

	if condition, ok := flattenedStep["condition"]; ok {
		step.Condition = deployments.DeploymentStepConditionType(condition.(string))
	}

	if conditionExpression, ok := flattenedStep["condition_expression"]; ok {
		step.Properties["Octopus.Step.ConditionVariableExpression"] = core.NewPropertyValue(conditionExpression.(string), false)
	}

	if packageRequirement, ok := flattenedStep["package_requirement"]; ok {
		step.PackageRequirement = deployments.DeploymentStepPackageRequirement(packageRequirement.(string))
	}

	if startTrigger, ok := flattenedStep["start_trigger"]; ok {
		step.StartTrigger = deployments.DeploymentStepStartTrigger(startTrigger.(string))
	}

	if targetRoles, ok := flattenedStep["target_roles"]; ok {
		step.Properties["Octopus.Action.TargetRoles"] = core.NewPropertyValue(strings.Join(getSliceFromTerraformTypeList(targetRoles), ","), false)
	}

	if windowSize, ok := flattenedStep["window_size"]; ok {
		step.Properties["Octopus.Action.MaxParallelism"] = core.NewPropertyValue(windowSize.(string), false)
	}

	var sort_order map[string]int = make(map[string]int)

	step_expansion := func(step_type_name string, step_type_action func(map[string]interface{}) *deployments.DeploymentAction) {
		if v, ok := flattenedStep[step_type_name]; ok {
			for _, tfAction := range v.([]interface{}) {
				flattenedAction := tfAction.(map[string]interface{})
				action := step_type_action(flattenedAction)
				step.Actions = append(step.Actions, action)

				// Pull out the sort_order if it exists. This is used to sort the actions later
				if posn, ok := flattenedAction["sort_order"].(int); ok && posn >= 0 {
					name := flattenedAction["name"].(string)
					sort_order[name] = posn
				}
			}
		}
	}

	step_expansion("action", expandAction)
	step_expansion("manual_intervention_action", expandManualInterventionAction)
	step_expansion("apply_terraform_template_action", expandApplyTerraformTemplateAction)
	step_expansion("deploy_package_action", expandDeployPackageAction)
	step_expansion("deploy_windows_service_action", expandDeployWindowsServiceAction)
	step_expansion("run_script_action", expandRunScriptAction)
	step_expansion("run_kubectl_script_action", expandRunKubectlScriptAction)
	step_expansion("deploy_kubernetes_secret_action", expandDeployKubernetesSecretAction)

	// Now that we have extracted all the steps off each of the properties into a single array, sort the array by the sort_order if provided
	if len(sort_order) > 0 {
		if len(sort_order) != len(step.Actions) {
			log.Printf("[WARN] Not all actions on step '%s' have a `sort_order` parameter so they may be sorted in an unexpected order", step.Name)
		}
		sort.SliceStable(step.Actions, func(i, j int) bool {
			return sort_order[step.Actions[i].Name] < sort_order[step.Actions[j].Name]
		})
	}

	return step
}

func flattenDeploymentSteps(deploymentSteps []*deployments.DeploymentStep) []map[string]interface{} {
	if deploymentSteps == nil {
		return nil
	}

	var flattenedDeploymentSteps = make([]map[string]interface{}, len(deploymentSteps))
	for key, deploymentStep := range deploymentSteps {
		flattenedDeploymentStep := map[string]interface{}{}
		flattenedDeploymentStep["condition"] = deploymentStep.Condition
		flattenedDeploymentStep["id"] = deploymentStep.ID
		flattenedDeploymentStep["name"] = deploymentStep.Name
		flattenedDeploymentStep["package_requirement"] = deploymentStep.PackageRequirement
		flattenedDeploymentStep["properties"] = flattenProperties(deploymentStep.Properties)
		flattenedDeploymentStep["start_trigger"] = deploymentStep.StartTrigger

		for propertyName, propertyValue := range deploymentStep.Properties {
			switch propertyName {
			case "Octopus.Action.TargetRoles":
				flattenedDeploymentStep["target_roles"] = strings.Split(propertyValue.Value, ",")
			case "Octopus.Action.MaxParallelism":
				flattenedDeploymentStep["window_size"] = propertyValue.Value
			case "Octopus.Step.ConditionVariableExpression":
				flattenedDeploymentStep["condition_expression"] = propertyValue.Value
			}
		}

		flatten_action_func := func(step_type_name string, i int, fp func(*deployments.DeploymentAction) map[string]interface{}) {

			if _, ok := flattenedDeploymentStep[step_type_name]; !ok {
				flattenedDeploymentStep[step_type_name] = make([]map[string]interface{}, 0)
			}

			action := fp(deploymentStep.Actions[i])
			action["sort_order"] = i + 1
			flattenedDeploymentStep[step_type_name] = append(flattenedDeploymentStep[step_type_name].([]map[string]interface{}), action)
		}

		for i := range deploymentStep.Actions {
			switch deploymentStep.Actions[i].ActionType {
			case "Octopus.KubernetesDeploySecret":
				flatten_action_func("deploy_kubernetes_secret_action", i, flattenDeployKubernetesSecretAction)
			case "Octopus.KubernetesRunScript":
				flatten_action_func("run_kubectl_script_action", i, flattenKubernetesRunScriptAction)
			case "Octopus.Manual":
				flatten_action_func("manual_intervention_action", i, flattenManualInterventionAction)
			case "Octopus.Script":
				flatten_action_func("run_script_action", i, flattenRunScriptAction)
			case "Octopus.TentaclePackage":
				flatten_action_func("deploy_package_action", i, flattenDeployPackageAction)
			case "Octopus.TerraformApply":
				flatten_action_func("apply_terraform_template_action", i, flattenApplyTerraformTemplateAction)
			case "Octopus.WindowsService":
				flatten_action_func("deploy_windows_service_action", i, flattenDeployWindowsServiceAction)
			default:
				flatten_action_func("action", i, flattenDeploymentAction)
			}

		}

		flattenedDeploymentSteps[key] = flattenedDeploymentStep
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
