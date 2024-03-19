package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func flattenGitDependency(gitDependency *gitdependencies.GitDependency) []interface{} {
	if gitDependency == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"repository_uri":      gitDependency.RepositoryUri,
		"default_branch":      gitDependency.DefaultBranch,
		"git_credential_type": gitDependency.GitCredentialType,
		"file_path_filters":   flattenArray(gitDependency.FilePathFilters),
		"git_credential_id":   gitDependency.GitCredentialId,
	}}
}

func expandGitDependency(set *schema.Set) *gitdependencies.GitDependency {
	if set == nil || len(set.List()) == 0 {
		return nil
	}

	flattenedMap := set.List()[0].(map[string]interface{})
	gitDependency := &gitdependencies.GitDependency{}

	if repositoryUri := flattenedMap["repository_uri"]; repositoryUri != nil {
		gitDependency.RepositoryUri = repositoryUri.(string)
	}

	if defaultBranch := flattenedMap["default_branch"]; defaultBranch != nil {
		gitDependency.DefaultBranch = defaultBranch.(string)
	}

	if gitCredentialType := flattenedMap["git_credential_type"]; gitCredentialType != nil {
		gitDependency.GitCredentialType = gitCredentialType.(string)
	}

	if filePathFilters := flattenedMap["file_path_filters"]; filePathFilters != nil {
		gitDependency.FilePathFilters = expandArray(filePathFilters.([]interface{}))
	}

	if gitCredentialId := flattenedMap["git_credential_id"]; gitCredentialId != nil {
		gitDependency.GitCredentialId = gitCredentialId.(string)
	}

	return gitDependency
}

func getGitDependencySchema(required bool) *schema.Schema {
	return &schema.Schema{
		Computed:    !required,
		Description: "Configuration for resource sourcing from a git repository.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"repository_uri": {
					Description:      "The Git URI for the repository where this resource is sourced from.",
					Required:         true,
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				},
				"default_branch": {
					Description:      "Name of the default branch of the repository.",
					Required:         true,
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				},
				"git_credential_type": {
					Description:      "The Git credential authentication type.",
					Required:         true,
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				},
				"file_path_filters": {
					Description: "List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.",
					Optional:    true,
					Type:        schema.TypeList,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"git_credential_id": {
					Description: "ID of an existing Git credential.",
					Optional:    true,
					Type:        schema.TypeString,
				},
			},
		},
		Optional: !required,
		Required: required,
		MaxItems: 1,
		Type:     schema.TypeSet,
	}
}

func addGitDependencySchema(element *schema.Resource) {
	element.Schema["git_dependency"] = getGitDependencySchema(false)
}
