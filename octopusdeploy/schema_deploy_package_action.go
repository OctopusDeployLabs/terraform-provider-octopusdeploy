package octopusdeploy

import (
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeployPackageAction(flattenedAction map[string]interface{}) *octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TentaclePackage"

	addWindowsServiceFeatureToActionResource(flattenedAction, action)
	return action
}

func flattenDeployPackageAction(action *octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.EnabledFeatures"]; ok {
		if strings.Contains(v.Value, "Octopus.Features.WindowsService") {
			flattenedAction["windows_service"] = flattenWindowsService(action.Properties)
		}
	}

	return flattenedAction
}

func getDeployPackageActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addPrimaryPackageSchema(element, true)
	// addCustomInstallationDirectoryFeature(element)
	// addIisWebSiteAndApplicationPoolFeature(element)
	addWindowsServiceFeature(element)
	// addCustomDeploymentScriptsFeature(element)
	// addJsonConfigurationVariablesFeature(element)
	// addConfigurationVariablesFeature(element)
	// addConfigurationTransformsFeature(element)
	// addSubstituteVariablesInFilesFeature(element)
	// addIis6HomeDirectoryFeature(element)
	// addRedGateDatabaseDeploymentFeature(element)
	return actionSchema
}
