package octopusdeploy

import (
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
	element.Schema[constScriptFileName] = &schema.Schema{
		Description: "The script file name in the package",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema[constScriptParameters] = &schema.Schema{
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func buildRunScriptActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := expandDeploymentAction(tfAction)
	resource.ActionType = "Octopus.Script"
	resource.Properties = merge(resource.Properties, buildRunScriptFromPackageActionResource(tfAction))

	variableSubstitutionInFiles := tfAction[constVariableSubstitutionInFiles].(string)
	if variableSubstitutionInFiles != "" {
		resource.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles
		resource.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"
		resource.Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
	}

	return resource
}

func buildRunScriptFromPackageActionResource(tfAction map[string]interface{}) map[string]string {
	return map[string]string{
		"Octopus.Action.Script.ScriptFileName":   tfAction["script_file_name"].(string),
		"Octopus.Action.Script.ScriptParameters": tfAction["script_parameters"].(string),
		"Octopus.Action.Script.ScriptSource":     "Package",
	}
}
