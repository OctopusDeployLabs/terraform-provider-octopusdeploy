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
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.Script"

	if v, ok := flattenedAction["script_body"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptBody"] = s
		}
	}

	if v, ok := flattenedAction["script_file_name"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptFileName"] = s
		}
	}

	if v, ok := flattenedAction["script_parameters"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptParameters"] = s
		}
	}

	if v, ok := flattenedAction["script_source"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.ScriptSource"] = s
		}
	}

	if v, ok := flattenedAction["script_syntax"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.Syntax"] = s
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
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptBody"]; ok {
		flattenedAction["script_body"] = v
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

	if v, ok := action.Properties["Octopus.Action.Script.Syntax"]; ok {
		flattenedAction["script_syntax"] = v
	}

	if v, ok := action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedAction["variable_substitution_in_files"] = v
	}

	return flattenedAction
}

func getRunScriptActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)

	element.Schema["script_body"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	element.Schema["script_syntax"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	element.Schema["variable_substitution_in_files"] = &schema.Schema{
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	return actionSchema
}
