package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGitTriggerSources(values interface{}) []filters.GitTriggerSource {
	if values == nil {
		return nil
	}

	var gitTriggerSources []filters.GitTriggerSource
	for _, v := range values.([]interface{}) {
		flattenedMap := v.(map[string]interface{})
		gitTriggerSources = append(gitTriggerSources, filters.GitTriggerSource{
			DeploymentActionSlug: flattenedMap["deployment_action_slug"].(string),
			GitDependencyName:    flattenedMap["git_dependency_name"].(string),
			IncludeFilePaths:     flattenedMap["include_file_paths"].([]string),
			ExcludeFilePaths:     flattenedMap["exclude_file_paths"].([]string),
		})
	}

	return gitTriggerSources
}

func flattenGitTriggerSources(gitTriggerSources []filters.GitTriggerSource) []interface{} {
	if len(gitTriggerSources) == 0 {
		return nil
	}

	flattenedGitTriggerSources := []interface{}{}
	for _, v := range gitTriggerSources {
		flattenedGitTriggerSources = append(flattenedGitTriggerSources, map[string]interface{}{
			"deployment_action_slug": v.DeploymentActionSlug,
			"git_dependency_name":    v.GitDependencyName,
			"include_file_paths":     v.IncludeFilePaths,
			"exclude_file_paths":     v.ExcludeFilePaths,
		})
	}

	return flattenedGitTriggerSources
}

func getGitTriggerSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deployment_action_slug": {
			Required: true,
			Type:     schema.TypeString,
		},
		"git_dependency_name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"include_file_paths": {
			Required: true,
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"exclude_file_paths": {
			Required: true,
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
