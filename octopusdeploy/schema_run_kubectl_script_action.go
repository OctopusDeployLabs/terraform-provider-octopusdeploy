package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getRunKubectlScriptSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)
	return actionSchema
}

func expandRunKubectlScriptAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	deploymentAction := expandDeploymentAction(flattenedAction)

	deploymentAction.ActionType = "Octopus.KubernetesRunScript"
	deploymentAction.Properties["Octopus.Action.Script.ScriptFileName"] = flattenedAction["script_file_name"].(string)
	deploymentAction.Properties["Octopus.Action.Script.ScriptParameters"] = flattenedAction["script_parameters"].(string)
	deploymentAction.Properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return deploymentAction
}

func flattenKubernetesRunScriptAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedKubernetesRunScriptAction := flattenCommonDeploymentAction(deploymentAction)

	if v, ok := deploymentAction.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v)
		flattenedKubernetesRunScriptAction["run_on_server"] = runOnServer
	}

	if v, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedKubernetesRunScriptAction["script_file_name"] = v
	}

	if v, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedKubernetesRunScriptAction["script_parameters"] = v
	}

	if v, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedKubernetesRunScriptAction["script_source"] = v
	}

	if v, ok := deploymentAction.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedKubernetesRunScriptAction["variable_substitution_in_files"] = v
	}

	return flattenedKubernetesRunScriptAction
}
