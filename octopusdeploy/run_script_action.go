package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getRunScriptActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)

	element.Schema["variable_substitution_in_files"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
	}

	return actionSchema
}

func addScriptFromPackageSchema(element *schema.Resource) {

	element.Schema["script_file_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The script file name in the package",
		Optional:    true,
	}

	element.Schema["script_parameters"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
	}
}

func buildRunScriptActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)

	resource.ActionType = "Octopus.Script"

	resource.Properties = merge(resource.Properties, buildRunScriptFromPackageActionResource(tfAction))

	variableSubstitutionInFiles := tfAction["variable_substitution_in_files"].(string)

	if variableSubstitutionInFiles != "" {
		resource.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles
		resource.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

		resource.Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
	}

	return resource
}

func buildRunScriptFromPackageActionResource(tfAction map[string]interface{}) map[string]string {

	properties := make(map[string]string)

	properties["Octopus.Action.Script.ScriptFileName"] = tfAction["script_file_name"].(string)
	properties["Octopus.Action.Script.ScriptParameters"] = tfAction["script_parameters"].(string)
	properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return properties
}
