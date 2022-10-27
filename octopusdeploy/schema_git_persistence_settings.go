package octopusdeploy

import (
	"context"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGitPersistenceSettings(values interface{}) projects.IPersistenceSettings {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	url, err := url.Parse(flattenedMap["url"].(string))
	if err != nil {
		return nil
	}

	var gitCredential credentials.IGitCredential
	if v, ok := flattenedMap["git_credential_id"]; ok {
		gitCredential = credentials.NewReference(v.(string))
	}

	if v, ok := flattenedMap["credentials"]; ok {
		gitCredential = expandGitCredential(v)
	} else {
		gitCredential = credentials.NewAnonymous()
	}

	return projects.NewGitPersistenceSettings(
		flattenedMap["base_path"].(string),
		nil,
		gitCredential,
		flattenedMap["default_branch"].(string),
		[]string{},
		url,
	)
}

func flattenGitPersistenceSettings(ctx context.Context, d *schema.ResourceData, persistenceSettings projects.IPersistenceSettings) []interface{} {
	if persistenceSettings == nil || persistenceSettings.GetType() == "Database" {
		return nil
	}

	gitPersistanceSettings := persistenceSettings.(*projects.GitPersistenceSettings)

	flattenedGitPersistenceSettings := make(map[string]interface{})
	flattenedGitPersistenceSettings["base_path"] = gitPersistanceSettings.BasePath
	flattenedGitPersistenceSettings["default_branch"] = gitPersistanceSettings.DefaultBranch

	switch gitPersistanceSettings.Credentials.GetType() {
	case credentials.GitCredentialTypeReference:
		referenceProjectGitCredential := gitPersistanceSettings.Credentials.(*credentials.Reference)
		flattenedGitPersistenceSettings["git_credential_id"] = referenceProjectGitCredential.Id
	case credentials.GitCredentialTypeUsernamePassword:
		flattenedGitPersistenceSettings["credentials"] = flattenGitCredential(ctx, d, gitPersistanceSettings.Credentials)
	}

	if gitPersistanceSettings.URL != nil {
		flattenedGitPersistenceSettings["url"] = gitPersistanceSettings.URL.String()
	}

	return []interface{}{flattenedGitPersistenceSettings}
}
