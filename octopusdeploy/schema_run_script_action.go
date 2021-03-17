package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func addScriptFromPackageSchema(element *schema.Resource) {
	element.Schema["script_file_name"] = &schema.Schema{
		Description: "The script file name in the package",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["script_parameters"] = &schema.Schema{
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["script_source"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}
}

func expandRunScriptAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandDeploymentAction(flattenedAction)
	action.ActionType = "Octopus.Script"
	action.Properties = merge(action.Properties, flattenRunScriptActionProperties(action))

	if v, ok := flattenedAction["run_on_server"]; ok {
		action.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(v.(bool))
	}

	if v, ok := flattenedAction["script_file_name"]; ok {
		if scriptFileName := v.(string); len(scriptFileName) > 0 {
			action.Properties["Octopus.Action.Script.ScriptFileName"] = scriptFileName
		}
	}

	if v, ok := flattenedAction["script_parameters"]; ok {
		if scriptParameters := v.(string); len(scriptParameters) > 0 {
			action.Properties["Octopus.Action.Script.ScriptParameters"] = scriptParameters
		}
	}

	if v, ok := flattenedAction["script_source"]; ok {
		if scriptSource := v.(string); len(scriptSource) > 0 {
			action.Properties["Octopus.Action.Script.ScriptSource"] = scriptSource
		}
	}

	if variableSubstitutionInFiles, ok := flattenedAction["variable_substitution_in_files"]; ok {
		action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles.(string)
		action.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

		if len(action.Properties["Octopus.Action.EnabledFeatures"]) == 0 {
			action.Properties["Octopus.Action.EnabledFeatures"] = "Octopus.Features.SubstituteInFiles"
		} else {
			action.Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
		}
	}

	return action
}

func flattenRunScriptAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenCommonDeploymentAction(action)

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

func flattenRunScriptActionProperties(action octopusdeploy.DeploymentAction) map[string]string {
	flattenedProperties := map[string]string{}

	if runOnServer, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		flattenedProperties["Octopus.Action.RunOnServer"] = runOnServer
	}

	if scriptFileName, ok := action.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedProperties["Octopus.Action.Script.ScriptFileName"] = scriptFileName
	}

	if scriptParameters, ok := action.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedProperties["Octopus.Action.Script.ScriptParameters"] = scriptParameters
	}

	if scriptSource, ok := action.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedProperties["Octopus.Action.Script.ScriptSource"] = scriptSource
	}

	return flattenedProperties
}

func getRunScriptActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)

	element.Schema["variable_substitution_in_files"] = &schema.Schema{
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	return actionSchema
}
