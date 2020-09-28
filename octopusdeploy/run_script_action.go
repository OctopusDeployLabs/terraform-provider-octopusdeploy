package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getRunScriptActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addScriptFromPackageSchema(element)
	addPackagesSchema(element, false)

	element.Schema[constVariableSubstitutionInFiles] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
	}

	return actionSchema
}

func addScriptFromPackageSchema(element *schema.Resource) {

	element.Schema[constScriptFileName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The script file name in the package",
		Optional:    true,
	}

	element.Schema[constScriptParameters] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
	}
}

func buildRunScriptActionResource(tfAction map[string]interface{}) model.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)

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

	properties := make(map[string]string)

	properties["Octopus.Action.Script.ScriptFileName"] = tfAction[constScriptFileName].(string)
	properties["Octopus.Action.Script.ScriptParameters"] = tfAction[constScriptParameters].(string)
	properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return properties
}
