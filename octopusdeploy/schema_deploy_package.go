package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func buildDeployPackageActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {

	action := buildDeploymentActionResource(tfAction)
	action.ActionType = "Octopus.TentaclePackage"

	if tfAction == nil {
		log.Println("Deploy Package Resource is nil. Please confirm the package resource")
	}

	addWindowsServiceFeatureToActionResource(tfAction, action)
	return action
}
