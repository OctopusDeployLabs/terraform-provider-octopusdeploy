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
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesRunScript"

	action.Properties["Octopus.Action.Script.ScriptFileName"] = flattenedAction["script_file_name"].(string)
	action.Properties["Octopus.Action.Script.ScriptParameters"] = flattenedAction["script_parameters"].(string)
	action.Properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return action
}

func flattenKubernetesRunScriptAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedAction["script_file_name"] = v
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedAction["script_parameters"] = v
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedAction["script_source"] = v
	}

	if v, ok := action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedAction["variable_substitution_in_files"] = v
	}

	return flattenedAction
}
