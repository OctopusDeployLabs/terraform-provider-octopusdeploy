package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getDeployPackageAction() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addPrimaryPackageSchema(element, true)
	//addCustomInstallationDirectoryFeature(element)
	//addIisWebSiteAndApplicationPoolFeature(element)
	addWindowsServiceFeature(element)
	//addCustomDeploymentScriptsFeature(element)
	//addJsonConfigurationVariablesFeature(element)
	//addConfigurationVariablesFeature(element)
	//addConfigurationTransformsFeature(element)
	//addSubstituteVariablesInFilesFeature(element)
	//addIis6HomeDirectoryFeature(element)
	//addRedGateDatabaseDeploymentFeature(element)
	return actionSchema
}

func buildDeployPackageActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := buildDeploymentActionResource(tfAction)
	action.ActionType = "Octopus.TentaclePackage"
	addWindowsServiceFeatureToActionResource(tfAction, action)
	return action
}
