package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentActionSlugPackages(values interface{}) []packages.DeploymentActionSlugPackage {
	if values == nil {
		return nil
	}

	actionPackages := []packages.DeploymentActionSlugPackage{}
	for _, v := range values.([]interface{}) {
		flattenedMap := v.(map[string]interface{})
		actionPackages = append(actionPackages, packages.DeploymentActionSlugPackage{
			DeploymentActionSlug: flattenedMap["deployment_action_slug"].(string),
			PackageReference:     flattenedMap["package_reference"].(string),
		})
	}
	return actionPackages
}

func expandDeploymentActionSlugPrimaryPackages(values interface{}) []packages.DeploymentActionSlugPackage {
	if values == nil {
		return nil
	}

	actionPackages := []packages.DeploymentActionSlugPackage{}
	for _, v := range values.([]interface{}) {
		flattenedMap := v.(map[string]interface{})
		actionPackages = append(actionPackages, packages.DeploymentActionSlugPackage{
			DeploymentActionSlug: flattenedMap["deployment_action_slug"].(string),
		})
	}
	return actionPackages
}

func flattenDeploymentActionSlugPackages(deploymentActionSlugPackages []packages.DeploymentActionSlugPackage) []interface{} {
	if len(deploymentActionSlugPackages) == 0 {
		return nil
	}

	flattenedDeploymentActionSlugPackages := []interface{}{}
	for _, v := range deploymentActionSlugPackages {
		if v.PackageReference != "" {
			flattenedDeploymentActionSlugPackage := map[string]interface{}{
				"deployment_action_slug": v.DeploymentActionSlug,
				"package_reference":      v.PackageReference,
			}
			flattenedDeploymentActionSlugPackages = append(flattenedDeploymentActionSlugPackages, flattenedDeploymentActionSlugPackage)
		}
	}
	return flattenedDeploymentActionSlugPackages
}

func flattenDeploymentActionSlugPrimaryPackages(deploymentActionSlugPackages []packages.DeploymentActionSlugPackage) []interface{} {
	if len(deploymentActionSlugPackages) == 0 {
		return nil
	}

	flattenedDeploymentActionSlugPackages := []interface{}{}
	for _, v := range deploymentActionSlugPackages {
		if v.PackageReference == "" {
			flattenedDeploymentActionSlugPackage := map[string]interface{}{
				"deployment_action_slug": v.DeploymentActionSlug,
			}
			flattenedDeploymentActionSlugPackages = append(flattenedDeploymentActionSlugPackages, flattenedDeploymentActionSlugPackage)
		}
	}
	return flattenedDeploymentActionSlugPackages
}

func getDeploymentActionSlugPackageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deployment_action_slug": {
			Required: true,
			Type:     schema.TypeString,
		},
		"package_reference": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}

func getDeploymentActionSlugPrimaryPackageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deployment_action_slug": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}
