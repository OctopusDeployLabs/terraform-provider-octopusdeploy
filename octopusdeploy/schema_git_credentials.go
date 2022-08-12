package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
)

func expandGitCredential(values interface{}) projects.IGitCredential {
	if values == nil {
		return projects.NewAnonymousGitCredential()
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return projects.NewAnonymousGitCredential()
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	return projects.NewUsernamePasswordGitCredential(
		flattenedMap["username"].(string),
		core.NewSensitiveValue(flattenedMap["password"].(string)),
	)
}

func flattenGitCredential(gitCredential projects.IGitCredential, password string) []interface{} {
	if gitCredential == nil {
		return nil
	}

	if gitCredential.GetType() == "UsernamePassword" {
		usernamePasswordCredential := gitCredential.(*projects.UsernamePasswordGitCredential)
		return []interface{}{map[string]interface{}{
			"password": password,
			"username": usernamePasswordCredential.Username,
		}}
	}

	return []interface{}{}
}
