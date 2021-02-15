package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeployPackageAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedWindowsService := map[string]interface{}{}
	flattenWindowsService(flattenedWindowsService, deploymentAction.Properties)

	return map[string]interface{}{
		"name":            deploymentAction.Name,
		"primary_package": flattenPackageReferences(deploymentAction.Packages),
		"windows_service": []interface{}{flattenedWindowsService},
	}
}

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

func expandDeployPackageAction(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	deploymentAction := expandDeploymentAction(tfAction)
	deploymentAction.ActionType = "Octopus.TentaclePackage"

	if tfAction == nil {
		log.Println("Deploy Package Resource is nil. Please confirm the package resource")
	}

	addWindowsServiceFeatureToActionResource(tfAction, deploymentAction)
	return deploymentAction
}
