package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandGitPersistenceSettings(values interface{}) octopusdeploy.IPersistenceSettings {
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

	var credential octopusdeploy.IGitCredential
	if v, ok := flattenedMap["credentials"]; ok {
		credential = expandGitCredential(v)
	} else {
		credential = octopusdeploy.NewAnonymousGitCredential()
	}

	return octopusdeploy.NewGitPersistenceSettings(
		flattenedMap["base_path"].(string),
		credential,
		flattenedMap["default_branch"].(string),
		url,
	)
}

func flattenGitPersistenceSettings(persistenceSettings octopusdeploy.IPersistenceSettings, password string) []interface{} {
	if persistenceSettings == nil {
		return nil
	}

	if persistenceSettings.GetType() == "Database" {
		return nil
	}

	gitPersistanceSettings := persistenceSettings.(*octopusdeploy.GitPersistenceSettings)

	flattenedGitPersistenceSettings := make(map[string]interface{})
	flattenedGitPersistenceSettings["base_path"] = gitPersistanceSettings.BasePath
	flattenedGitPersistenceSettings["credentials"] = flattenGitCredential(gitPersistanceSettings.Credentials, password)
	flattenedGitPersistenceSettings["default_branch"] = gitPersistanceSettings.DefaultBranch

	if gitPersistanceSettings.URL != nil {
		flattenedGitPersistenceSettings["url"] = gitPersistanceSettings.URL.String()
	}

	return []interface{}{flattenedGitPersistenceSettings}
}
