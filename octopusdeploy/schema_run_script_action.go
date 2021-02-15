package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func expandRunScriptAction(flattenedRunScriptAction map[string]interface{}) octopusdeploy.DeploymentAction {
	deploymentAction := expandDeploymentAction(flattenedRunScriptAction)
	deploymentAction.ActionType = "Octopus.Script"
	deploymentAction.Properties = merge(deploymentAction.Properties, flattenRunScriptActionProperties(deploymentAction))

	if v, ok := flattenedRunScriptAction["run_on_server"]; ok {
		deploymentAction.Properties["Octopus.Action.RunOnServer"] = strconv.FormatBool(v.(bool))
	}

	if v, ok := flattenedRunScriptAction["script_file_name"]; ok {
		if scriptFileName := v.(string); len(scriptFileName) > 0 {
			deploymentAction.Properties["Octopus.Action.Script.ScriptFileName"] = scriptFileName
		}
	}

	if v, ok := flattenedRunScriptAction["script_parameters"]; ok {
		if scriptParameters := v.(string); len(scriptParameters) > 0 {
			deploymentAction.Properties["Octopus.Action.Script.ScriptParameters"] = scriptParameters
		}
	}

	if v, ok := flattenedRunScriptAction["script_source"]; ok {
		if scriptSource := v.(string); len(scriptSource) > 0 {
			deploymentAction.Properties["Octopus.Action.Script.ScriptSource"] = scriptSource
		}
	}

	if variableSubstitutionInFiles, ok := flattenedRunScriptAction["variable_substitution_in_files"]; ok {
		deploymentAction.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles.(string)
		deploymentAction.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

		if len(deploymentAction.Properties["Octopus.Action.EnabledFeatures"]) == 0 {
			deploymentAction.Properties["Octopus.Action.EnabledFeatures"] = "Octopus.Features.SubstituteInFiles"
		} else {
			deploymentAction.Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
		}
	}

	return deploymentAction
}

func flattenRunScriptAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	return map[string]interface{}{
		"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
		"channels":                           deploymentAction.Channels,
		"container":                          flattenDeploymentActionContainer(deploymentAction.Container),
		"condition":                          deploymentAction.Condition,
		"environments":                       deploymentAction.Environments,
		"excluded_environments":              deploymentAction.ExcludedEnvironments,
		"id":                                 deploymentAction.ID,
		"is_disabled":                        deploymentAction.IsDisabled,
		"is_required":                        deploymentAction.IsRequired,
		"name":                               deploymentAction.Name,
		"notes":                              deploymentAction.Notes,
		"package":                            flattenPackageReferences(deploymentAction.Packages),
		"properties":                         flattenRunScriptActionProperties(deploymentAction),
		"tenant_tags":                        deploymentAction.TenantTags,
	}
}

func flattenRunScriptActionProperties(deploymentAction octopusdeploy.DeploymentAction) map[string]string {
	flattenedRunScriptAction := map[string]string{}

	if runOnServer, ok := deploymentAction.Properties["Octopus.Action.RunOnServer"]; ok {
		flattenedRunScriptAction["run_on_server"] = runOnServer
	}

	if scriptFileName, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptFileName"]; ok {
		flattenedRunScriptAction["script_file_name"] = scriptFileName
	}

	if scriptParameters, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptParameters"]; ok {
		flattenedRunScriptAction["script_parameters"] = scriptParameters
	}

	if scriptSource, ok := deploymentAction.Properties["Octopus.Action.Script.ScriptSource"]; ok {
		flattenedRunScriptAction["script_source"] = scriptSource
	}

	return flattenedRunScriptAction
}
