package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandVersionControlSettings(values interface{}) *octopusdeploy.VersionControlSettings {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	if flattenedMap["type"] == "Database" {
		return &octopusdeploy.VersionControlSettings{
			Type: "Database",
		}
	}

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

	return octopusdeploy.NewVersionControlSettings(
		flattenedMap["base_path"].(string),
		credential,
		flattenedMap["default_branch"].(string),
		"VersionControlled",
		url,
	)
}
