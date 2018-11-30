package octopusdeploy

import (
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getRunRunKubectlScriptSchema()  *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addPrimaryPackageSchema(element, false)
	addScriptFromPackageSchema(element)

	return actionSchema;
}

func buildRunKubectlScriptActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)

	resource.ActionType = "Octopus.KubernetesRunScript"

	resource.Properties = merge(resource.Properties, buildScriptFromPackageProperties(tfAction))

	resource.Properties["Octopus.Action.Script.ScriptFileName"] = tfAction["script_file_name"].(string)
	resource.Properties["Octopus.Action.Script.ScriptParameters"] = tfAction["script_parameters"].(string)
	resource.Properties["Octopus.Action.Script.ScriptSource"] = "Package"


	return resource
}
func merge(map1 map[string]string, map2 map[string]string) map[string]string {

}
func buildScriptFromPackageProperties(tfAction map[string]interface{}) map[string]string {
	properties := make(map[string]string)
	properties["Octopus.Action.Script.ScriptFileName"] = tfAction["script_file_name"].(string)
	properties["Octopus.Action.Script.ScriptParameters"] = tfAction["script_parameters"].(string)
	properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return properties;
}
