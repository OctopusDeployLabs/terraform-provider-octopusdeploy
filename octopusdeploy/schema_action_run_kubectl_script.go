package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getRunKubectlScriptSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)
	return actionSchema
}

func expandRunKubectlScriptAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesRunScript"

	action.Properties["Octopus.Action.Script.ScriptFileName"] = core.NewPropertyValue(flattenedAction["script_file_name"].(string), false)
	action.Properties["Octopus.Action.Script.ScriptParameters"] = core.NewPropertyValue(flattenedAction["script_parameters"].(string), false)
	action.Properties["Octopus.Action.Script.ScriptSource"] = core.NewPropertyValue("Package", false)

	return action
}

func flattenKubernetesRunScriptAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedAction["script_file_name"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedAction["script_parameters"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedAction["script_source"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedAction["variable_substitution_in_files"] = v.Value
	}

	return flattenedAction
}
