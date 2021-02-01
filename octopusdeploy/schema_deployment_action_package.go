package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentActionPackage(values interface{}) *octopusdeploy.DeploymentActionPackage {
	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.DeploymentActionPackage{
		DeploymentAction: flattenedMap["deployment_action"].(string),
		PackageReference: flattenedMap["package_reference"].(string),
	}
}

func flattenDeploymentActionPackage(deploymentActionPackage *octopusdeploy.DeploymentActionPackage) []interface{} {
	if deploymentActionPackage == nil {
		return nil
	}

	flattenedDeploymentActionPackage := make(map[string]interface{})
	flattenedDeploymentActionPackage["deployment_action"] = deploymentActionPackage.DeploymentAction
	flattenedDeploymentActionPackage["package_reference"] = deploymentActionPackage.PackageReference
	return []interface{}{flattenedDeploymentActionPackage}
}

func getDeploymentActionPackageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deployment_action": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"package_reference": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
