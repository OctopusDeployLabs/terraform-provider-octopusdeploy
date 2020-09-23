package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func getDeployPackageAction() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
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

func buildDeployPackageActionResource(tfAction map[string]interface{}) model.DeploymentAction {

	action := buildDeploymentActionResource(tfAction)
	action.ActionType = "Octopus.TentaclePackage"
	addWindowsServiceFeatureToActionResource(tfAction, action)
	return action
}
