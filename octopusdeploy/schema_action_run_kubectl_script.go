package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getRunKubectlScriptSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	element.Schema["script_body"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addWorkerPoolSchema(element)
	addWorkerPoolVariableSchema(element)
	addPackagesSchema(element, false)
	return actionSchema
}

func expandRunKubectlScriptAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesRunScript"

	if v, ok := flattenedAction["script_body"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptBody"] = core.NewPropertyValue(s, false)
		}
	}

	if v, ok := flattenedAction["script_file_name"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptFileName"] = core.NewPropertyValue(s, false)
		}
	}

	if v, ok := flattenedAction["script_parameters"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptParameters"] = core.NewPropertyValue(s, false)
		}
	}

	if v, ok := flattenedAction["script_source"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptSource"] = core.NewPropertyValue(s, false)
		}
	}
	return action
}

func flattenKubernetesRunScriptAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedAction["run_on_server"] = runOnServer
	}

	if len(action.WorkerPool) > 0 {
		flattenedAction["worker_pool_id"] = action.WorkerPool
	}

	if len(action.WorkerPoolVariable) > 0 {
		flattenedAction["worker_pool_variable"] = action.WorkerPoolVariable
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedAction["script_file_name"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedAction["script_parameters"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptBody"]; ok {
		flattenedAction["script_body"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedAction["script_source"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedAction["variable_substitution_in_files"] = v.Value
	}

	return flattenedAction
}
