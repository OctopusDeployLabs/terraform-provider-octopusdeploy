package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandGitCredential(values interface{}) octopusdeploy.IGitCredential {
	if values == nil {
		return octopusdeploy.NewAnonymousGitCredential()
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return octopusdeploy.NewAnonymousGitCredential()
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	return octopusdeploy.NewUsernamePasswordGitCredential(
		flattenedMap["username"].(string),
		octopusdeploy.NewSensitiveValue(flattenedMap["password"].(string)),
	)
}

func flattenGitCredential(gitCredential octopusdeploy.IGitCredential, password string) []interface{} {
	if gitCredential == nil {
		return nil
	}

	if gitCredential.GetType() == "UsernamePassword" {
		usernamePasswordCredential := gitCredential.(*octopusdeploy.UsernamePasswordGitCredential)
		return []interface{}{map[string]interface{}{
			"password": password,
			"username": usernamePasswordCredential.Username,
		}}
	}

	return []interface{}{}
}
