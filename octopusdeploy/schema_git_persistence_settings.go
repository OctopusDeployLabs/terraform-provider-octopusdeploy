package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
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

	var credential projects.IGitCredential
	if v, ok := flattenedMap["credentials"]; ok {
		credential = expandGitCredential(v)
	} else {
		credential = projects.NewAnonymousGitCredential()
	}

	return projects.NewGitPersistenceSettings(
		flattenedMap["base_path"].(string),
		credential,
		flattenedMap["default_branch"].(string),
		url,
	)
}

func flattenGitPersistenceSettings(persistenceSettings projects.IPersistenceSettings, password string) []interface{} {
	if persistenceSettings == nil {
		return nil
	}

	if persistenceSettings.GetType() == "Database" {
		return nil
	}

	gitPersistanceSettings := persistenceSettings.(*projects.GitPersistenceSettings)

	flattenedGitPersistenceSettings := make(map[string]interface{})
	flattenedGitPersistenceSettings["base_path"] = gitPersistanceSettings.BasePath
	flattenedGitPersistenceSettings["credentials"] = flattenGitCredential(gitPersistanceSettings.Credentials, password)
	flattenedGitPersistenceSettings["default_branch"] = gitPersistanceSettings.DefaultBranch

	if gitPersistanceSettings.URL != nil {
		flattenedGitPersistenceSettings["url"] = gitPersistanceSettings.URL.String()
	}

	return []interface{}{flattenedGitPersistenceSettings}
}
