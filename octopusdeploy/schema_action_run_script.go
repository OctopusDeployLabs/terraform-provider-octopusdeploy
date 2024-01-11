package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
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

func expandRunScriptAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	if len(flattenedAction) == 0 {
		return nil
	}

	action := expandAction(flattenedAction)
	if action == nil {
		return nil
	}

	action.ActionType = "Octopus.Script"

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

	if v, ok := flattenedAction["script_syntax"]; ok {
		if s := v.(string); len(s) > 0 {
			action.Properties["Octopus.Action.Script.Syntax"] = core.NewPropertyValue(s, false)
		}
	}

	if variableSubstitutionInFiles, ok := flattenedAction["variable_substitution_in_files"]; ok {
		action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = core.NewPropertyValue(variableSubstitutionInFiles.(string), false)
		action.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = core.NewPropertyValue("True", false)

		if len(action.Properties["Octopus.Action.EnabledFeatures"].Value) == 0 {
			action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue("Octopus.Features.SubstituteInFiles", false)
		} else {
			actionProperty := action.Properties["Octopus.Action.EnabledFeatures"].Value + ",Octopus.Features.SubstituteInFiles"
			action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue(actionProperty, false)
		}
	}

	if v, ok := flattenedAction["worker_pool_id"]; ok {
		action.WorkerPool = v.(string)
	}

	if v, ok := flattenedAction["worker_pool_variable"]; ok {
		action.WorkerPoolVariable = v.(string)
	}

	return action
}

func flattenRunScriptAction(action *deployments.DeploymentAction) map[string]interface{} {
	if action == nil {
		return nil
	}

	flattenedAction := flattenAction(action)

	if len(action.WorkerPool) > 0 {
		flattenedAction["worker_pool_id"] = action.WorkerPool
	}

	if len(action.WorkerPoolVariable) > 0 {
		flattenedAction["worker_pool_variable"] = action.WorkerPoolVariable
	}

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.Script.ScriptBody"]; ok {
		flattenedAction["script_body"] = v.Value
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

	if v, ok := action.Properties["Octopus.Action.Script.Syntax"]; ok {
		flattenedAction["script_syntax"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"]; ok {
		flattenedAction["variable_substitution_in_files"] = v.Value
	}

	return flattenedAction
}

func getRunScriptActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addPropertiesSchema(element, "This attribute is deprecated and will be removed in a future release. Please use the attributes that match the properties that are stored to this map.")
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)
	addWorkerPoolSchema(element)
	addWorkerPoolVariableSchema(element)

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
