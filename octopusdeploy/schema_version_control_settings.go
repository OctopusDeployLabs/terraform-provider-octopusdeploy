package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
)

func expandVersionControlSettings(values interface{}) *projects.VersionControlSettings {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	if flattenedMap["type"] == "Database" {
		return &projects.VersionControlSettings{
			Type: "Database",
		}
	}

	url, err := url.Parse(flattenedMap["url"].(string))
	if err != nil {
		return nil
	}

	var credential credentials.IGitCredential
	if v, ok := flattenedMap["credentials"]; ok {
		credential = expandGitCredential(v)
	} else {
		credential = credentials.NewAnonymous()
	}

	return projects.NewVersionControlSettings(
		flattenedMap["base_path"].(string),
		credential,
		flattenedMap["default_branch"].(string),
		"VersionControlled",
		url,
	)
}
