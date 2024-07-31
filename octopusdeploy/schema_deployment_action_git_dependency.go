package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentActionGitDependencies(deploymentActionGitDependencies []gitdependencies.DeploymentActionGitDependency) []interface{} {
	if len(deploymentActionGitDependencies) == 0 {
		return nil
	}

	var flattenedDeploymentActionPackages []interface{}
	for _, v := range deploymentActionGitDependencies {
		flattenedDeploymentActionPackage := map[string]interface{}{
			"deployment_action_slug": v.DeploymentActionSlug,
			"git_dependency_name":    v.GitDependencyName,
		}
		flattenedDeploymentActionPackages = append(flattenedDeploymentActionPackages, flattenedDeploymentActionPackage)
	}
	return flattenedDeploymentActionPackages
}

func expandDeploymentActionGitDependencies(values interface{}) []gitdependencies.DeploymentActionGitDependency {
	if values == nil {
		return nil
	}

	var gitDependencies []gitdependencies.DeploymentActionGitDependency
	for _, v := range values.([]interface{}) {
		flattenedMap := v.(map[string]interface{})
		gitDependencies = append(gitDependencies, gitdependencies.DeploymentActionGitDependency{
			DeploymentActionSlug: flattenedMap["deployment_action_slug"].(string),
			GitDependencyName:    flattenedMap["git_dependency_name"].(string),
		})
	}
	return gitDependencies
}

func getDeploymentActionGitDependencySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deployment_action_slug": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"git_dependency_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
